package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"

	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/commands/common/config"
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
	"setup-auto-complete",
	"update",
	"config",
}

// isHelpAlias returns true when a command's arguments are equivalent to the
// help command. For example, `lingo review --help` == `lingo help review`.
func isHelpAlias(flags []string) bool {
	helpFlags := strings.Split(cli.HelpFlag.Name, ", ")
	return len(flags) == 1 && (flags[0] == "--"+helpFlags[0] || flags[0] == "-"+helpFlags[1])
}

func needsLingoHome(cmd string, flags []string) bool {
	if isHelpAlias(flags) {
		return false
	}
	for _, c := range cmdNeedsLingoHome {
		if c == cmd {
			return true
		}
	}
	return false
}

func needsDotLingo(cmd string, flags []string) bool {
	if isHelpAlias(flags) {
		return false
	}
	for _, c := range cmdNeedsDotLingo {
		if c == cmd {
			return true
		}
	}
	return false
}

func BeforeCMD(c *cli.Context) error {
	var currentCMDName string
	var flags []string
	args := c.Args()
	if args.Present() {
		currentCMDName = args.First()
		flags = args.Tail()
	}

	if needsLingoHome(currentCMDName, flags) {
		ensureLingoHome()
	}

	// ensure we have a .lingo file
	if needsDotLingo(currentCMDName, flags) {
		if cfgPath, _ := common.TenetCfgPath(c); cfgPath == "" {
			return errors.Wrap(common.ErrMissingDotLingo, errors.New("ui"))
		}
	}

	return nil
}

func ensureLingoHome() error {
	// create lingo home if it doesn't exist
	lHome := util.MustLingoHome()
	if _, err := os.Stat(lHome); os.IsNotExist(err) {
		err := os.MkdirAll(lHome, 0775)
		if err != nil {
			return errors.Trace(err)
		}
	}

	tenetsHome := filepath.Join(lHome, "tenets")
	if _, err := os.Stat(tenetsHome); os.IsNotExist(err) {
		err := os.MkdirAll(tenetsHome, 0775)
		if err != nil {
			return errors.Trace(err)
		}
	}

	logsHome := filepath.Join(lHome, "logs")
	if _, err := os.Stat(logsHome); os.IsNotExist(err) {
		err := os.MkdirAll(logsHome, 0775)
		if err != nil {
			return errors.Trace(err)
		}
	}

	scriptsHome := filepath.Join(lHome, "scripts")
	if _, err := os.Stat(scriptsHome); os.IsNotExist(err) {
		err := os.MkdirAll(scriptsHome, 0775)
		if err != nil {
			fmt.Printf("WARNING: could not create scripts directory: %v \n", err)
		}
	}

	ensureConfigs(lHome)

	return nil
}

func ensureConfigs(lHome string) {

	configsHome := filepath.Join(lHome, "configs")
	if _, err := os.Stat(configsHome); os.IsNotExist(err) {
		err := os.MkdirAll(configsHome, 0775)
		if err != nil {
			fmt.Printf("WARNING: could not create configs directory: %v \n", err)
		}
	}

	servicesCfg := filepath.Join(configsHome, config.ServicesCfgFile)
	if _, err := os.Stat(servicesCfg); os.IsNotExist(err) {
		err := ioutil.WriteFile(servicesCfg, []byte(config.ServicesTmpl), 0644)
		if err != nil {
			fmt.Printf("WARNING: could not create services config: %v \n", err)
		}
	}

	defaultsCfg := filepath.Join(configsHome, config.DefaultsCfgFile)
	if _, err := os.Stat(defaultsCfg); os.IsNotExist(err) {
		err := ioutil.WriteFile(defaultsCfg, []byte(config.DefaultsTmpl), 0644)
		if err != nil {
			fmt.Printf("WARNING: could not create services config: %v \n", err)
		}
	}
}
