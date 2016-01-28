package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"

	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/util"
)

var InfoCMD = cli.Command{
	Name:  "info",
	Usage: "show information about a tenet",
	Description: `
	"lingo info <tenet-name>"
`[1:],

	Action: infoAction,
	BashComplete: func(c *cli.Context) {
		// This will complete if no args are passed
		if len(c.Args()) > 0 {
			return
		}

		tenets, err := util.BinTenets()
		if err != nil {
			log.Printf("auto complete error %v", err)
			return
		}

		for _, t := range tenets {
			fmt.Println(t)
		}
	},
}

func infoAction(ctx *cli.Context) {
	var commandIsTenet bool
	var cfg common.TenetConfig
	if a := ctx.Args(); len(a) != 1 {
		common.OSErrf(" info expects one argument, the tenet name")
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
			common.OSErrf("tenet not found")
			return
		}
	}

	tnCMDs, err := newTenetCMDs(ctx, cfg)
	if err != nil {
		common.OSErrf(err.Error())
		return
	}
	defer tnCMDs.closeService()

	if err := tnCMDs.printInfo(); err != nil {
		common.OSErrf(err.Error())
		return
	}
}
