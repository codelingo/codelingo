package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/util"
)

var InitCMD = cli.Command{
	Name:   "init",
	Usage:  "create a " + defaultTenetCfgPath + " config file in current or specified directory",
	Action: initLingo,
}

var configSeed = `
# This config file is in TOML format. Manually add and configure tenets
# following these guidelines: http://lingo.reviews/getting-started#the-lingo-
# file. Alternatively, use the add/remove commands on the lingo app.

`[1:]

// TODO(waigani) set lingo-home flag and test init creates correct home dir.

func initLingo(c *cli.Context) {
	if err := maxArgs(c, 1); err != nil {
		oserrf(err.Error())
		return
	}

	// create lingo home if it doesn't exist
	lHome := util.MustLingoHome()
	if _, err := os.Stat(lHome); os.IsNotExist(err) {
		err := os.MkdirAll(lHome, 0644)
		if err != nil {
			panic(err)
		}
	}

	tenetsHome := filepath.Join(lHome, "tenets")
	if _, err := os.Stat(tenetsHome); os.IsNotExist(err) {
		err := os.MkdirAll(tenetsHome, 0644)
		if err != nil {
			panic(err)
		}
	}

	// Create a new tenet config file at either the provided directory or
	// a location from flags or environment, or the current directory
	cfgPath, _ := filepath.Abs(desiredTenetCfgPath(c))
	if len(c.Args()) > 0 {
		cfgPath, _ = filepath.Abs(c.Args()[0])

		// Check that it exists and is a directory
		if pathInfo, err := os.Stat(cfgPath); os.IsNotExist(err) {
			oserrf("directory %q not found", cfgPath)
		} else if !pathInfo.IsDir() {
			oserrf("%q is not a directory", cfgPath)
		}

		// Use default config filename
		cfgPath = filepath.Join(cfgPath, defaultTenetCfgPath)
	}

	if _, err := os.Stat(cfgPath); err == nil {
		oserrf("Already initialised using tenet config file %q", cfgPath)
	}

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	err := enc.Encode(&config{Include: "*", Cascade: true})
	if err != nil {
		oserrf(err.Error())
		return
	}

	if err = ioutil.WriteFile(cfgPath, buf.Bytes(), 0644); err != nil {
		oserrf(err.Error())
		return
	}

	fmt.Printf("Successfully initialised. Lingo config file written to %q\n", cfgPath)
}
