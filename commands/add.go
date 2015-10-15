package commands

import (
	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/tenet"
)

var AddCMD = cli.Command{
	Name:  "add",
	Usage: "add a tenet to lingo",
	Description: `

  "lingo remove github.com/lingo-reviews/unused-args"

`[1:],
	Action: add,
}

func add(c *cli.Context) {
	if l := len(c.Args()); l != 1 {
		oserrf("expected 1 argument, got %d", l)
		return
	}
	cfg, err := readTenetCfgFile(c)
	if err != nil {
		oserrf("reading config file: %s", err.Error())
		return
	}

	imageName := c.Args().First()

	if hasTenet(cfg, imageName) {
		oserrf(`error: tenet "%s" already added`, imageName)
		return
	}

	cfg.Configs = append(cfg.Configs, tenet.Config{Name: imageName})

	if err := writeTenetCfgFile(c, cfg); err != nil {
		oserrf(err.Error())
		return
	}

	// TODO(waigani) open an interactive shell, prompt to set options.
}
