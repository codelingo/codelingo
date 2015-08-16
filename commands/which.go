package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var WhichCMD = cli.Command{
	Name:        "which",
	Usage:       "prints path to tenet.toml",
	Description: "prints path to tenet.toml",
	Action:      which,
}

func which(c *cli.Context) {
	path, err := tenetCfgPath(c)
	if err != nil {
		oserrf(err.Error())
	}
	fmt.Println(path)
}
