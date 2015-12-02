package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/util"
)

// List of commands which can be run without needing a config file
// TODO(matt) remove synonimous commands from here and use a resolving function
// once UI for that is sorted out

// A list of cmds that need a .lingo file
var cmdNeedsDotLingo = []string{
	"add",
	"remove",
	"rm",
	"review",
	"pull",
	"list",
	"ls",
	"write-docs",
	"docs",
	"edit",
}

var cmdNeedsLingoHome = []string{
	"build",
	"init",
	"add",
	"remove",
	"rm",
	"review",
	"pull",
	"list",
	"ls",
	"write-docs",
	"docs",
	"edit",
}

func needsLingoHome(cmd string) bool {
	for _, c := range cmdNeedsLingoHome {
		if c == cmd {
			return true
		}
	}
	return false
}

func needsDotLingo(cmd string) bool {
	for _, c := range cmdNeedsDotLingo {
		if c == cmd {
			return true
		}
	}
	return false
}

func BeforeCMD(c *cli.Context) error {
	var currentCMDName string
	args := c.Args()
	if args.Present() {
		currentCMDName = args.First()
	}

	if needsLingoHome(currentCMDName) {
		ensureLingoHome()
	}

	// ensure we have a .lingo file
	if needsDotLingo(currentCMDName) {
		if cfgPath, _ := tenetCfgPath(c); cfgPath == "" {
			return errors.Wrap(errors.New("No .lingo configuration found. Run `lingo init` to create a .lingo file in the current directory"), errors.New("ui"))
		}
	}

	return nil
}

func ensureLingoHome() {
	// create lingo home if it doesn't exist
	lHome := util.MustLingoHome()
	if _, err := os.Stat(lHome); os.IsNotExist(err) {
		err := os.MkdirAll(lHome, 0775)
		if err != nil {
			panic(err)
		}
	}

	tenetsHome := filepath.Join(lHome, "tenets")
	if _, err := os.Stat(tenetsHome); os.IsNotExist(err) {
		err := os.MkdirAll(tenetsHome, 0775)
		if err != nil {
			panic(err)
		}
	}

	logsHome := filepath.Join(lHome, "logs")
	if _, err := os.Stat(logsHome); os.IsNotExist(err) {
		err := os.MkdirAll(logsHome, 0775)
		if err != nil {
			panic(err)
		}
	}

	scriptsHome := filepath.Join(lHome, "scripts")
	if _, err := os.Stat(scriptsHome); os.IsNotExist(err) {
		err := os.MkdirAll(scriptsHome, 0775)
		if err != nil {
			fmt.Printf("WARNING: could not create scripts directory: %v \n", err)
		}
	}

}
