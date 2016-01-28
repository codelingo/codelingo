package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"

	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/util"
)

var RemoveCMD = cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Usage:   "remove a tenet from lingo",
	Description: `

  "lingo remove github.com/lingo-reviews/unused-args"

`[1:],
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "group",
			Value: "default",
			Usage: "group to remove tenet from"},
	},
	Action: remove,
	BashComplete: func(c *cli.Context) {
		// This will complete if no args are passed
		if len(c.Args()) > 0 {
			return
		}

		// TODO(waigani) read from .lingo not bin tenets
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

func remove(c *cli.Context) {
	if l := len(c.Args()); l != 1 {
		common.OSErrf("error: expected 1 argument, got %d", l)
		return
	}

	cfgPath, err := common.TenetCfgPath(c)
	if err != nil {
		common.OSErrf("reading config file: %s", err.Error())
		return
	}
	cfg, err := common.BuildConfig(cfgPath, common.CascadeNone)
	if err != nil {
		common.OSErrf("reading config file: %s", err.Error())
		return
	}

	imageName := c.Args().First()

	if !cfg.HasTenet(imageName) {
		common.OSErrf(`error: tenet "%s" not found in %q`, imageName, c.GlobalString(common.TenetCfgFlg.Long))
		return
	}

	if err := cfg.RemoveTenet(imageName, c.String("group")); err != nil {
		common.OSErrf(err.Error())
		return
	}

	if err := common.WriteConfigFile(c, cfg); err != nil {
		common.OSErrf(err.Error())
		return
	}
}
