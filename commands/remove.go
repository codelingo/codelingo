package commands

import (
	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/tenet"
)

var RemoveCMD = cli.Command{Name: "remove",
	Usage: "remove a tenet from lingo",
	Description: `

  "lingo remove github.com/lingo-reviews/unused-args"

`[1:],
	Action: remove,
}

func remove(c *cli.Context) {
	if l := len(c.Args()); l != 1 {
		oserrf("error: expected 1 argument, got %d", l)
	}

	cfg, err := readTenetCfgFile(c)
	if err != nil {
		oserrf("reading config file: %s", err.Error())
	}

	imageName := c.Args().First()

	if !hasTenet(cfg, imageName) {
		oserrf(`error: tenet "%s" not found in %q`, imageName, c.GlobalString(tenetCfgFlg.long))
	}

	var tenets []tenet.Tenet
	for _, t := range cfg.Tenets {
		if t.Name != imageName {
			tenets = append(tenets, t)
		}
	}
	cfg.Tenets = tenets

	if err := writeTenetCfgFile(c, cfg); err != nil {
		oserrf(err.Error())
	}
}
