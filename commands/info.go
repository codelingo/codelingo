package commands

import "github.com/codegangsta/cli"

var InfoCMD = cli.Command{
	Name:  "info",
	Usage: "show information about a tenet",
	Description: `
	"lingo info <tenet-name>"
`[1:],

	Action: infoAction,
}

func infoAction(ctx *cli.Context) {
	var commandIsTenet bool
	var cfg TenetConfig
	if a := ctx.Args(); len(a) != 1 {
		oserrf(" info expects one argument, the tenet name")
		return
	} else {
		// Does the command match an installed tenet?
		for _, cfg = range listTenets(ctx) {
			if a[0] == cfg.Name {
				commandIsTenet = true
				break
			}
		}
		if !commandIsTenet {
			oserrf("tenet not found")
			return
		}
	}

	tnCMDs, err := newTenetCMDs(ctx, cfg)
	if err != nil {
		oserrf(err.Error())
		return
	}
	defer tnCMDs.closeService()

	if err := tnCMDs.printInfo(); err != nil {
		oserrf(err.Error())
		return
	}
}
