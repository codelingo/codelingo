package commands

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/util"
)

var ListCMD = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "list tenets",
	Description: "Lists all tenets added to .lingo run.",
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
		cmd := exec.Command("tree", filepath.Join(util.MustLingoHome(), "tenets"))
		cmd.Stdout = &stdout
		// cmd.Stderr = &stderr
		cmd.Run()

		fmt.Print(string(stdout.Bytes()))
		return
	}

	for _, t := range listTenets(c) {
		fmt.Printf("%s (%s)\n", t.Name, t.Driver)
	}

}

func listTenets(c *cli.Context) []TenetConfig {
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
	return cfg.AllTenets()
}
