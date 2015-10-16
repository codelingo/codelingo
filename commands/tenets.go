package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var TenetsCMD = cli.Command{
	Name:        "tenets",
	Usage:       "list tenets",
	Description: "Lists all tenets added to .lingo, run `lingo help <tenet-name>` to see options. Default options are set in .lingo",
	Action:      tenetsAction,
}

func tenetsAction(c *cli.Context) {
	for _, t := range tenetCfgs(c) {
		fmt.Println(t.Name)
	}
}
