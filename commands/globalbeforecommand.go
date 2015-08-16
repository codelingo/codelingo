package commands

import (
	"github.com/codegangsta/cli"
	"github.com/juju/errors"
)

func BeforeCMD(c *cli.Context) error {
	var currentCMDName string
	args := c.Args()
	if args.Present() {
		currentCMDName = args.First()
	}

	// ensure we have a tenet.toml file
	if currentCMDName != "init" {
		if cfgPath, _ := tenetCfgPath(c); cfgPath == "" {
			desiredCfg := desiredTenetCfgPath(c)
			return errors.Errorf("not a lingo project. %s not found (nor in any of the parent directories). Run `lingo init` to write a tenet.toml file in the current directory", desiredCfg)
		}
	}

	return nil
}
