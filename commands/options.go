package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var OptionsCMD = cli.Command{
	Name:        "options",
	Usage:       "options <tenet-url>",
	Description: "configure tenet options",
	Action:      options,
}

func options(c *cli.Context) {
	fmt.Print("cfg options not implemented")
}
