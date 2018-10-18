package rewrite

import (
	"strings"

	"github.com/codegangsta/cli"

	rewriterpc "github.com/codelingo/codelingo/flows/codelingo/rewrite/rpc"
	flowutil "github.com/codelingo/codelingo/sdk/flow"
	"github.com/fatih/color"
	"github.com/golang/protobuf/proto"
)

var DecoratorCMD = &flowutil.DecoratorCommand{
	Command: cli.Command{
		Name:  "rewrite",
		Usage: "Modify code following tenets in codelingo.yaml.",
		Flags: []cli.Flag{

			cli.BoolFlag{
				Name:  "replace",
				Usage: "replace the decorated node with the new source code",
			},
			cli.BoolFlag{
				Name:  "prepend",
				Usage: "prepend the decorated node with the new source code",
			},
			cli.BoolFlag{
				Name:  "append",
				Usage: "append the decorated node with the new source code",
			},
			cli.BoolFlag{
				Name:  "line",
				Usage: "operate on the entire line instead of the byte offsets",
			},
			cli.BoolFlag{
				Name:  "byte",
				Usage: "operate on the byte offsets instead of the entire line",
			},
			cli.BoolFlag{
				Name:  "start-offset",
				Usage: "operate on the start offset only",
			},
			cli.BoolFlag{
				Name:  "end-offset",
				Usage: "operate on the end offset only",
			},
			cli.BoolFlag{
				Name:  "start-to-end-offset",
				Usage: "operate on the start to end offset range",
			},
		},
		Description: `
"@rewrite rewrites the decorated node.
`[1:],
	},
	ConfirmDecorated: decoratorAction,
}

func decoratorAction(ctx *cli.Context, payload proto.Message) (bool, error) {

	hunk := payload.(*rewriterpc.Hunk)

	item := &flowutil.ConfirmerItem{
		Preview: func() string {

			g := color.New(color.FgGreen).SprintfFunc()
			return indent(g("\n%s", hunk.SRC), true, false)
		},

		// Options is a map of option keys to confirm functions. Each confirm function returns: <keep>bool, <retry>bool, <err>error
		Options: map[string]func() (bool, bool, error){
			"[o]pen": func() (bool, bool, error) {
				return flowutil.OpenFileConfirmAction(hunk.Filename, int64(0))
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
			"[k]eep":    []string{"k", "K", "keep"},
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
