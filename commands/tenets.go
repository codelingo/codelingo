package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/tenet"
)

var TenetsCMD = cli.Command{
	Name:        "tenets",
	Usage:       "list tenets",
	Description: "Lists all tenets added to .lingo, run `lingo help <tenet-name>` to see options. Default options are set in .lingo",
	Action:      tenetsAction,
}

func tenetsAction(c *cli.Context) {
	for _, t := range tenets(c) {
		fmt.Println(t.String())
	}
}

func tenets(c *cli.Context) []tenet.Tenet {
	cfg, err := readTenetCfgFile(c)
	if err != nil {
		oserrf("reading config file: %s", err.Error())
		return nil
	}
	return cfg.Tenets
}
