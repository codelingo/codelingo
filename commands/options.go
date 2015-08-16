package commands

import (
	"github.com/codegangsta/cli"
	"github.com/waigani/xxx"
)

var OptionsCMD = cli.Command{
	Name:        "options",
	Usage:       "options <tenet-url>",
	Description: "configure tenet options",
	Action:      options,
}

func options(c *cli.Context) {
	xxx.Print("cfg options")
}
