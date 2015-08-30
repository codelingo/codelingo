package review

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/lingo-reviews/dev/tenet"
	"github.com/skratchdot/open-golang/open"
	"github.com/waigani/diffparser"
)

type issueConfirmer struct {
	confidence  float64
	userConfirm bool
	inDiff      func(*tenet.Issue) bool
	// TODO(waigani) make this a func var instead
	hostAbsBasePath string
}

func NewConfirmer(c *cli.Context) *issueConfirmer {
	cfm := issueConfirmer{
		confidence:      c.Float64("confidence"),
		userConfirm:     !c.Bool("no-user-confirm"),
		hostAbsBasePath: hostAbsBasePath(c),
	}

	if c.Bool("diff") {
		cfm.inDiff = newInDiffFunc()
	}

	return &cfm
}

func hostAbsBasePath(c *cli.Context) string {
	p := c.GlobalString("repo-path")
	basePath, err := filepath.Abs(p)
	if err != nil {
		// TODO(waigani) return err and have NewConfirmer return err also.
		panic(err)
	}
	return basePath
}

// TODO(waigani) screen diff tenet side - see other diff comment.
func newInDiffFunc() func(issue *tenet.Issue) bool {
	diff := diffparser.Parse(rawDiff())
	diffChanges := diff.Changed()

	return func(issue *tenet.Issue) bool {
		start := issue.Position.Start.Line
		end := start
		if issue.Position.End != nil {
			end = issue.Position.End.Line
		}

		// was issue found in diff?
		if newLines, ok := diffChanges[issue.Filename()]; ok {
			for _, lineNo := range newLines {
				if lineNo >= start && lineNo <= end {
					return true
				}
			}
		}
		return false
	}
}

// TODO(waigani) this just reads unstaged changes from git in pwd. Change diff
// from a flag to a sub command which pipes args to git diff.
func rawDiff() string {
	var stdout bytes.Buffer
	c := exec.Command("git", "diff")
	c.Stdout = &stdout
	// c.Stderr = &stderr
	c.Run()

	return string(stdout.Bytes())
}

// TODO(waigani) make this configurable.
// understandsLines is a list of apps that understand line number prepended to
// a filename.
var understandsLines = map[string]bool{
	"subl":    true,
	"sublime": true,
}

var editor string

// confirm returns true if the issue should be kept or false if it should be
// dropped.
func (c issueConfirmer) Confirm(attempt int, issue *tenet.Issue) bool {
	if attempt == 0 {
		if issue.Confidence < c.confidence ||
			(c.inDiff != nil && !c.inDiff(issue)) {
			return false
		}
		if !c.userConfirm {
			return true
		}
	}
	if attempt == 0 {
		fmt.Println(FormatPlainText(issue))
	}

	attempt++
	var options string
	fmt.Print("\n[o]pen [d]iscard [K]eep:")
	fmt.Scanln(&options)
	// if err != nil || n != 1 {
	// 	// TODO(waigani)  handle invalid input
	// 	fmt.Println("invalid input", n, err)
	// }
	switch options {
	case "o":
		var app string
		defaultEditor := "optional"
		if editor != "" {
			defaultEditor = editor
		}
		fmt.Printf("application (%s):", defaultEditor)
		fmt.Scanln(&app)
		filename := c.hostFilePath(issue.Filename())
		if app != "" {
			if _, ok := understandsLines[app]; ok {
				filename += fmt.Sprintf(":%d", issue.Position.Start.Line)
			}
			err := open.StartWith(filename, app)
			if err != nil {
				fmt.Println(err)
			}
			editor = app
		} else {
			var err error
			if defaultEditor == "optional" {
				err = open.Start(filename)
			} else {
				err = open.StartWith(filename, defaultEditor)
			}
			if err != nil {
				fmt.Println(err)
			}
		}
		c.Confirm(attempt, issue)
	case "d":
		return false
	case "", "k", "K", "\n":
		return true
	default:
		fmt.Printf("invalid input: %s\n", options)
		fmt.Println(options)
		c.Confirm(attempt, issue)
	}

	return true
}

func (c *issueConfirmer) hostFilePath(file string) string {
	return strings.Replace(file, "/source", c.hostAbsBasePath, 1)
}

// TODO(waigani) remove dependency on dev/tenet. Use a simpler internal
// representation of tenet.Issue.
func FormatPlainText(issue *tenet.Issue) string {
	m := color.New(color.FgWhite, color.Faint).SprintfFunc()
	y := color.New(color.FgYellow).SprintfFunc()
	yf := color.New(color.FgYellow, color.Faint).SprintfFunc()
	c := color.New(color.FgCyan).SprintfFunc()

	address := m("%s-%d:%d", issue.Position.Start.String(), issue.Position.End.Line, issue.Position.End.Column)
	comment := strings.Trim(issue.Comment, "\n")
	comment = c(indent("\n"+comment+"\n", false))

	ctxBefore := indent(yf("\n...\n%s", issue.CtxBefore), false)
	issueLines := indent(y("\n%s", issue.LineText), true)
	ctxAfter := indent(yf("\n%s\n...", issue.CtxAfter), false)
	src := ctxBefore + issueLines + ctxAfter

	return fmt.Sprintf("%s\n%s\n%s", address, comment, src)
}

func indent(str string, point bool) string {
	replace := "\n    "
	if point {
		replace = "\n  > "
	}
	return strings.Replace(str, "\n", replace, -1)
}
