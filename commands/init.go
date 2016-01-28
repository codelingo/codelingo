package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/commands/common"
)

var InitCMD = cli.Command{
	Name:   "init",
	Usage:  "create a " + common.DefaultTenetCfgPath + " config file in the current directory",
	Action: initLingo,
}

var configSeed = `
# This config file is in TOML format. Manually add and configure tenets
# following these guidelines: http://lingo.reviews/getting-started#the-lingo-
# file. Alternatively, use the add/remove commands on the lingo app.

`[1:]

// TODO(waigani) set lingo-home flag and test init creates correct home dir.

func initLingo(c *cli.Context) {
	if err := common.MaxArgs(c, 1); err != nil {
		common.OSErrf(err.Error())
		return
	}

	// Create a new tenet config file at either the provided directory or
	// a location from flags or environment, or the current directory
	cfgPath, _ := filepath.Abs(common.DesiredTenetCfgPath(c))
	if len(c.Args()) > 0 {
		cfgPath, _ = filepath.Abs(c.Args()[0])

		// Check that it exists and is a directory
		if pathInfo, err := os.Stat(cfgPath); os.IsNotExist(err) {
			common.OSErrf("directory %q not found", cfgPath)
		} else if !pathInfo.IsDir() {
			common.OSErrf("%q is not a directory", cfgPath)
		}

		// Use default config filename
		cfgPath = filepath.Join(cfgPath, common.DefaultTenetCfgPath)
	}

	if _, err := os.Stat(cfgPath); err == nil {
		common.OSErrf("Already initialised using tenet config file %q", cfgPath)
	}

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	err := enc.Encode(&common.Config{Include: "*", Cascade: true})
	if err != nil {
		common.OSErrf(err.Error())
		return
	}

	if err = ioutil.WriteFile(cfgPath, buf.Bytes(), 0644); err != nil {
		common.OSErrf(err.Error())
		return
	}
	fmt.Printf("Successfully initialised. Lingo config file written to %q\n", cfgPath)
}
