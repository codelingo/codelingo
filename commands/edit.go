package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/util"
)

var EditCMD = cli.Command{
	Name:  "edit",
	Usage: "edit the .lingo file",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "editor, e",
			Value:  "vi",
			Usage:  "editor to open config with",
			EnvVar: "LINGO_EDITOR",
		},
	},
	Action: edit,
}

func edit(c *cli.Context) {
	cfg, err := common.TenetCfgPath(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd, err := util.OpenFileCmd(c.String("editor"), cfg, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd.Run()
}
