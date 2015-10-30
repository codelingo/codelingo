package commands

import (
	"fmt"
	"os/exec"

	"github.com/codegangsta/cli"
)

var EditCMD = cli.Command{
	Name:  "edit",
	Usage: "edit the configuration file",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "editor, e",
			Value:  "subl", // TODO(waigani) DEMOWARE
			Usage:  "editor to open config with",
			EnvVar: "LINGO_EDITOR",
		},
	},
	Action: edit,
}

func edit(c *cli.Context) {
	cfg, err := tenetCfgPath(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO(waigani) DEMOWARE
	cmd := exec.Command(c.String("editor"), cfg)
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	cmd.Run()
}
