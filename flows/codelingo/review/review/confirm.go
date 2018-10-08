package review

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/rpc/flow"
	"github.com/fatih/color"
	"github.com/waigani/diffparser"
)

type IssueConfirmer struct {
	keepAll bool
	output  bool
}

func NewConfirmer(output, keepAll bool, d *diffparser.Diff) (*IssueConfirmer, error) {
	cfm := IssueConfirmer{
		keepAll: keepAll,
		output:  output,
	}

	return &cfm, nil
}

func GetDiffRootPath(filename string) string {
	// Get filename relative to git root folder
	// TODO: Handle error in case of git not being installed
	// https://github.com/codelingo/demo/issues/2
	out, err := exec.Command("git", "ls-tree", "--full-name", "--name-only", "HEAD", filename).Output()
	if err == nil && len(out) != 0 {
		if len(out) != 0 {
			filename = strings.TrimSuffix(string(out), "\n")
		}
	}
	return filename
}

var editor string

// confirm returns true if the issue should be kept or false if it should be
// dropped.
func (c IssueConfirmer) Confirm(attempt int, issue *flow.Issue) bool {
	if c.keepAll {
		return true
	}
	if attempt == 0 {
		fmt.Println(c.FormatPlainText(issue))
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
		// https://github.com/codelingo/demo/issues/3
		if editor != "" {
			defaultEditor = editor
		}
		fmt.Printf("application (%s):", defaultEditor)
		fmt.Scanln(&app)
		filename := issue.Position.Start.Filename
		if app == "" {
			app = defaultEditor
		}
		// c := issue.Position.Start.Column // TODO(waigani) use column
		// https://github.com/codelingo/demo/issues/4
		l := issue.Position.Start.Line
		cmd, err := util.OpenFileCmd(app, filename, l)
		if err != nil {
			fmt.Println(err)
			return c.Confirm(attempt, issue)
		}

		if err = cmd.Start(); err != nil {
			log.Println(err)
		}
		if err = cmd.Wait(); err != nil {
			log.Println(err)
		}

		editor = app

		c.Confirm(attempt, issue)
	case "d":
		issue.Discard = true

		// TODO(waigani) only prompt for reason if we're sending to a service.
		// https://github.com/codelingo/demo/issues/5
		//fmt.Print("reason: ")
		//in := bufio.NewReader(os.Stdin)
		//issue.DiscardReason, _ = in.ReadString('\n')

		// TODO(waigani) we are now always returning true. Need to decide
		// how caller will deal with removing isseus, ie. KeptIssues vs AllIssues,
		// then being returning false here
		// https://github.com/codelingo/demo/issues/6
		return true
	case "", "k", "K", "\n":
		return true
	default:
		fmt.Printf("invalid input: %s\n", options)
		fmt.Println(options)
		c.Confirm(attempt, issue)
	}

	return true
}

func (c *IssueConfirmer) FormatPlainText(issue *flow.Issue) string {
	m := color.New(color.FgWhite, color.Faint).SprintfFunc()
	y := color.New(color.FgYellow).SprintfFunc()
	yf := color.New(color.FgYellow, color.Faint).SprintfFunc()
	cy := color.New(color.FgCyan).SprintfFunc()
	filename := issue.Position.Start.Filename

	addrFmtStr := fmt.Sprintf("%s:%d", filename, issue.Position.End.Line)
	if col := issue.Position.End.Column; col != 0 {
		addrFmtStr += fmt.Sprintf(":%d", col)
	}
	address := m(addrFmtStr)
	comment := strings.Trim(issue.Comment, "\n")
	comment = cy(indent("\n"+comment+"\n", false))

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
