package main

import (
	"github.com/codegangsta/cli"
	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/sdk/flow"
	"github.com/juju/errors"
)

var docsCommand = cli.Command{
	Name:  "docs",
	Usage: "Generate documentation from Tenets",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  util.OutputFlg.String(),
			Usage: "File to save found results to.",
		},
		cli.StringFlag{
			Name:  "template, t",
			Value: "default",
			Usage: "The template to use when generating docs.",
		},
	},
	Description: `
""$ lingo docs" .
`[1:],
	Action: docsAction,
}

func main() {
	if err := flow.Run(docsCommand); err != nil {
		flow.HandleErr(err)
	}
}

func docsAction(ctx *cli.Context) {
	docs, err := docsCMD(ctx)
	if err != nil {

		// Debugging
		util.Logger.Debugw("docsAction", "err_stack", errors.ErrorStack(err))

		util.FatalOSErr(err)
		return
	}

	if ctx.IsSet("template") || ctx.IsSet("t") {
		print("template")
	} else {

		print(docs)
	}

}

func docsCMD(cliCtx *cli.Context) (string, error) {
	return "docs", nil
}
