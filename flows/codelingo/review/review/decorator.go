package review

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"

	flowutil "github.com/codelingo/codelingo/sdk/flow"
	"github.com/codelingo/rpc/flow"
	"github.com/fatih/color"
	"github.com/golang/protobuf/proto"
)

var DecoratorApp = &flowutil.DecoratorApp{
	App: cli.App{
		Name:  "review",
		Usage: "Comment on the decorated code.",
		Flags: []cli.Flag{},
	},
	ConfirmDecorated: decoratorAction,
	SetUserVar: func(v *flowutil.UserVar) {
		v.SetAsDefault()
	},
	// help info
	DecoratorUsage:   "<comment>",
	DecoratorExample: `"this is a review comment"`,
}

func decoratorAction(ctx *cli.Context, payload proto.Message) (bool, error) {

	issue := payload.(*flow.Issue)

	item := &flowutil.ConfirmerItem{
		Preview: func() string {
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
			comment = cy(indent("\n"+comment+"\n", false, false))

			ctxBefore := indent(yf("\n...\n%s", issue.CtxBefore), false, false)
			issueLines := indent(y("\n%s", issue.LineText), true, false)
			ctxAfter := indent(yf("\n%s\n...", issue.CtxAfter), false, false)
			src := ctxBefore + issueLines + ctxAfter

			return fmt.Sprintf("%s\n%s\n%s", address, comment, src)
		},

		// Options is a map of option keys to confirm functions. Each confirm function returns: <keep>bool, <retry>bool, <err>error
		Options: map[string]func() (bool, bool, error){
			"[o]pen": func() (bool, bool, error) {
				return flowutil.OpenFileConfirmAction(issue.Position.Start.Filename, int64(0))
			},
			"[k]eep": func() (bool, bool, error) {
				return true, false, nil
			},
			"[d]iscard": func() (bool, bool, error) {
				return false, false, nil
			},
		},

		// OptionKeyMap maps an option key to aliases e.g. "[k]eep" => "k", "keep", "K"
		OptionKeyMap: map[string][]string{
			"[o]pen":    []string{"o", "O", "open"},
			"[k]eep":    []string{"k", "K", "keep", ""},
			"[d]iscard": []string{"d", "D", "Discard"},
		},

		// FlagOptions is a map of flags to options show if that flag is present.
		// "_all" is a builtin flag for those options that don't require a flag.
		FlagOptions: map[string][]string{
			"_all": []string{"[o]pen", "[k]eep", "[d]iscard"},
		},
	}

	return item.Confirm(ctx)
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
