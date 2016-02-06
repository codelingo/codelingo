package commands

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/commands/common/config"
	"github.com/lingo-reviews/lingo/util"

	"github.com/codegangsta/cli"
)

var ConfigCMD = cli.Command{
	Name:  "config",
	Usage: "open lingo's configuration files",
	Description: `

  $ lingo config <config-file>

  Avaliable configs are: "defaults" and "services"

  If you have auto-completion enabled, tab after "$ lingo config" to see the list of available configs.

`[1:],
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "editor",
			Value: "vi",
			Usage: "editor to open the config file with",
		},
	},
	Action: editConfigAction,
	BashComplete: func(c *cli.Context) {
		// This will complete if no args are passed
		if len(c.Args()) > 0 {
			return
		}

		cfgDir, err := util.ConfigHome()
		if err != nil {
			log.Printf("auto complete error %v", err)
			return
		}

		files, err := filepath.Glob(cfgDir + "/*")
		if err != nil {
			log.Printf("auto complete error %v", err)
			return
		}

		cfgs := make([]string, len(files))
		for i, f := range files {
			f = strings.TrimPrefix(f, cfgDir+"/")
			f = strings.TrimSuffix(f, ".yaml")
			cfgs[i] = f
		}

		for _, c := range cfgs {
			fmt.Println(c)
		}
	},
}

func editConfigAction(ctx *cli.Context) {
	cfg := config.DefaultsCfgFile
	if args := ctx.Args(); len(args) > 0 {
		cfg = args[0] + ".yaml"
	}

	if err := config.Edit(cfg, "vi"); err != nil {
		common.OSErrf(err.Error())
	}
}
