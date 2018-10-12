package flow

// TODO(waigani) currently hardcoded to rewrite. Generalise this so it can be
// used for all Flow endpoints. It may be possible to use the command config
// structs to do this.

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	rewriterpc "github.com/codelingo/codelingo/flows/codelingo/rewrite/rpc"

	"github.com/briandowns/spinner"
	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/rpc/flow"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/waigani/diffparser"
	"github.com/waigani/xxx"
)

type hunkconfirmer struct {
	keepAll bool
	output  bool
}

func NewConfirmer(output, keepAll bool, d *diffparser.Diff) (*hunkconfirmer, error) {
	cfm := hunkconfirmer{
		keepAll: keepAll,
		output:  output,
	}

	return &cfm, nil
}

func ConfirmIssues(cancel context.CancelFunc, hunkc chan *rewriterpc.Hunk, errorc chan error, keepAll bool, saveToFile string) ([]*rewriterpc.Hunk, error) {
	defer util.Logger.Sync()

	var confirmedHunks []*rewriterpc.Hunk
	spnr := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spnr.Start()
	defer spnr.Stop()

	output := saveToFile == ""
	cfm, err := NewConfirmer(output, keepAll, nil)

	if err != nil {
		cancel()
		return nil, errors.Trace(err)
	}

	// If user is manually confirming reviews, set a long timeout.
	timeout := time.After(time.Hour * 1)
	if keepAll {
		timeout = time.After(time.Minute * 11)
	}

l:
	for {
		select {
		case err, ok := <-errorc:
			if !ok {
				errorc = nil
				break
			}

			// Abort review
			cancel()
			util.Logger.Debugf("error: %s", errors.ErrorStack(err))
			return nil, errors.Trace(err)
		case hunk, ok := <-hunkc:
			if !keepAll {
				spnr.Stop()
			}
			if !ok {
				hunkc = nil
				break
			}

			if cfm.Confirm(0, hunk) {
				confirmedHunks = append(confirmedHunks, hunk)
			}

			if !keepAll {
				spnr.Restart()
			}
		case <-timeout:
			cancel()
			return nil, errors.New("timed out waiting for response")
		}
		if hunkc == nil && errorc == nil {
			break l
		}
	}

	// Stop spinner if it hasn't been stopped already
	if keepAll {
		spnr.Stop()
	}
	return confirmedHunks, nil
}

// returns the full lines of the SRC for the hunk.
func fullLineSRC(hunk *flow.Issue, newSRC string) (string, error) {

	pos := hunk.GetPosition()
	startPos := pos.GetStart()
	endPos := pos.GetEnd()

	file, err := os.Open(startPos.Filename)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var byt int64
	var src []byte
	var strtLineByt int64
	var endLineByt int64
	var foundStart bool
	var foundEnd bool
	for scanner.Scan() {
		src = append(src, append(scanner.Bytes(), []byte("\n")...)...)

		startByt := byt
		endByt := startByt + int64(len(scanner.Bytes())+1) // +1 for \n char

		if startPos.Offset >= startByt && startPos.Offset <= endByt {

			// found start line
			strtLineByt = startByt
			foundStart = true
		}

		if endPos.Offset >= startByt && endPos.Offset <= endByt {

			// found end line
			endLineByt = endByt
			foundEnd = true
		}

		if foundStart && foundEnd {

			// xxx.Print(string(src))
			// fmt.Printf("\n[%d:%d]", strtLineByt, startPos.Offset)
			beforeNewSRC := string(src[strtLineByt:startPos.Offset])
			endNewSRC := string(src[endPos.Offset:endLineByt])

			return beforeNewSRC + newSRC + endNewSRC, nil

		}

		byt = endByt
	}

	return "", errors.Trace(scanner.Err())
}

func GetDiffRootPath(filename string) string {
	// Get filename relative to git root folder
	// TODO: Handle error in case of git not being installed
	// https://github.com/codelingo/demo/hunks/2
	out, err := exec.Command("git", "ls-tree", "--full-name", "--name-only", "HEAD", filename).Output()
	if err == nil && len(out) != 0 {
		if len(out) != 0 {
			filename = strings.TrimSuffix(string(out), "\n")
		}
	}
	return filename
}

var editor string

// confirm returns true if the hunk should be kept or false if it should be
// dropped.
func (c hunkconfirmer) Confirm(attempt int, hunk *rewriterpc.Hunk) bool {
	if c.keepAll {
		return true
	}
	if attempt == 0 {
		fmt.Println(c.FormatPlainText(hunk))
	}
	attempt++
	var options string
	fmt.Print("\n[o]pen")
	if c.output {
		fmt.Print(" [d]iscard [k]eep")
	}
	fmt.Print(": ")

	fmt.Scanln(&options)

	switch options {
	case "o":
		var app string
		defaultEditor := "vi" // TODO(waigani) use EDITOR or VISUAL env vars
		// https://github.com/codelingo/demo/hunks/3
		if editor != "" {
			defaultEditor = editor
		}
		fmt.Printf("application (%s):", defaultEditor)
		fmt.Scanln(&app)
		filename := hunk.Filename
		if app == "" {
			app = defaultEditor
		}
		// c := hunk.Position.Start.Column // TODO(waigani) use column
		// https://github.com/codelingo/demo/hunks/4

		// TODO(waigani) calc line from offset
		l := int64(0)
		cmd, err := util.OpenFileCmd(app, filename, l)
		if err != nil {
			fmt.Println(err)
			return c.Confirm(attempt, hunk)
		}

		if err = cmd.Start(); err != nil {
			log.Println(err)
		}
		if err = cmd.Wait(); err != nil {
			log.Println(err)
		}

		editor = app

		c.Confirm(attempt, hunk)
	case "d":
		return false
	case "", "k", "K", "\n":
		return true
	default:
		fmt.Printf("invalid input: %s\n", options)
		fmt.Println(options)
		c.Confirm(attempt, hunk)
	}

	// TODO(waigani) build up hunks here.

	return true
}

func (c *hunkconfirmer) FormatPlainText(hunk *rewriterpc.Hunk) string {

	xxx.Dump(hunk)

	g := color.New(color.FgGreen).SprintfFunc()
	return indent(g("\n%s", hunk.SRC), true, false)

	// TODO(waigani) generate a diff hunk

	// m := color.New(color.FgWhite, color.Faint).SprintfFunc()
	// y := color.New(color.FgRed).SprintfFunc()
	// yf := color.New(color.FgWhite, color.Faint).SprintfFunc()
	// filename := hunk.Filename

	// // TODO(waigani) get line from offset
	// line := 0
	// addrFmtStr := fmt.Sprintf("%s:%d", filename, line)

	// // TODO(waigani) get column from offset
	// col := 0
	// addrFmtStr += fmt.Sprintf(":%d", col)
	// address := m(addrFmtStr)

	// ctxBefore := indent(yf("\n...\n%s", hunk.CtxBefore), false, false)
	// oldLines := indent(y("\n%s", hunk.LineText), false, true)

	// newLines := indent(g("\n%s", newSRC), true, false)
	// ctxAfter := indent(yf("\n%s\n...", hunk.CtxAfter), false, false)
	// src := ctxBefore + oldLines + newLines + ctxAfter

	// return fmt.Sprintf("%s\n%s", address, src)
}

func indent(str string, add, remove bool) string {
	replace := "\n    "
	if add {
		replace = "\n  + "
	}
	if remove {
		replace = "\n  - "
	}
	return strings.Replace(str, "\n", replace, -1)
}
