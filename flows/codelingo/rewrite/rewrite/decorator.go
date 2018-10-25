package rewrite

import (
	"strings"

	"github.com/urfave/cli"

	rewriterpc "github.com/codelingo/codelingo/flows/codelingo/rewrite/rpc"
	flowutil "github.com/codelingo/codelingo/sdk/flow"
	"github.com/fatih/color"
	"github.com/golang/protobuf/proto"
)

var DecoratorApp = &flowutil.DecoratorApp{
	App: cli.App{
		Name:  "rewrite",
		Usage: "Replace the decorated node with the new code",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "replace, r",
				Usage: "replace the decorated node with the new source code",
			},
			cli.BoolFlag{
				Name:  "prepend, p",
				Usage: "prepend the decorated node with the new source code",
			},
			cli.BoolFlag{
				Name:  "append, a",
				Usage: "append the decorated node with the new source code",
			},
			cli.BoolFlag{
				Name:  "line, l",
				Usage: "operate on the entire line instead of the byte offsets",
			},
			cli.BoolFlag{
				Name:  "byte, b",
				Usage: "operate on the byte offsets instead of the entire line",
			},
			cli.BoolFlag{
				Name:  "start-offset, s",
				Usage: "operate on the start offset only",
			},
			cli.BoolFlag{
				Name:  "end-offset, e",
				Usage: "operate on the end offset only",
			},
			cli.BoolFlag{
				Name:  "start-to-end-offset, m",
				Usage: "operate on the start to end offset range",
			},
			cli.StringFlag{
				Name:  "new-file, n",
				Usage: "create a new file",
			},
			cli.IntFlag{
				Name:  "new-file-perm, x",
				Usage: "file permission for new file",
			},
		},
	},
	ConfirmDecorated: decoratorAction,
	DecoratorUsage:   "[options] <new code>",
	DecoratorExample: `--prepend --line "// new comment on a Golang function`,
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
