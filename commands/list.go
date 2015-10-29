package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/tenet"
)

var ListCMD = cli.Command{
	Name:        "list",
	Usage:       "list tenets",
	Description: "Lists all tenets added to tenet.toml run.",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "installed",
			Usage: "list all tenets installed on this machine",
		},
	},
	Action: listAction,
}

func listAction(c *cli.Context) {

	// TODO(waigani) DEMOWARE
	if c.Bool("installed") {
		var stdout bytes.Buffer
		cmd := exec.Command("tree", path.Join(os.Getenv("HOME"), ".lingo/tenets/"))
		cmd.Stdout = &stdout
		// cmd.Stderr = &stderr
		cmd.Run()

		fmt.Print(string(stdout.Bytes()))
		return
	}

	for _, t := range listTenets(c) {
		fmt.Println(t.Name)
	}

}

func listTenets(c *cli.Context) []tenet.Config {
	cfgPath, err := tenetCfgPath(c)
	if err != nil {
		oserrf("could not read configuration: %s", err.Error())
		return nil
	}
	cfg, err := buildConfig(cfgPath, CascadeNone)
	if err != nil {
		oserrf("could not read configuration: %s", err.Error())
		return nil
	}
	return cfg.Tenets
}
