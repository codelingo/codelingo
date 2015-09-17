package commands

import (
	"github.com/codegangsta/cli"
	"github.com/juju/errors"
)

// List of commands which can be run without needing a config file
// TODO(matt) remove synonimous commands from here and use a resolving function
// once UI for that is sorted out
var standaloneCommands = []string{"", "init", "help", "h"}

func BeforeCMD(c *cli.Context) error {
	var currentCMDName string
	args := c.Args()
	if args.Present() {
		currentCMDName = args.First()
	}

	// ensure we have a tenet.toml file
	standalone := false
	for _, c := range standaloneCommands {
		if c == currentCMDName {
			standalone = true
			break
		}
	}
	if !standalone {
		if cfgPath, _ := tenetCfgPath(c); cfgPath == "" {
			return errors.Wrap(errors.New("No tenet.toml configuration found. Run `lingo init` to create a tenet.toml file in the current directory"), errors.New("ui"))
		}
	}

	return nil
}
