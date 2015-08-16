package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/codegangsta/cli"
)

var InitCMD = cli.Command{
	Name:  "init",
	Usage: "Create a tenet.toml config file in pwd",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "pwd",
			Usage: "initialise the current directory to use it's own set of tenets",
		},
	},
	Action: initLingo,
}

var configSeed = `
# This config file is in TOML format. Manually add and configure tenets
# following these guidelines: http://lingo.reviews/getting-started#the-lingo-
# file. Alternatively, use the add/remove commands on the lingo app.

`[1:]

// TODO(waigani) set lingo-home flag and test init creates correct home dir.

func initLingo(c *cli.Context) {
	initPWD := c.Bool("pwd")
	// first init lingo
	home := c.GlobalString(lingoHomeFlg.long)
	// CONTINUE HERE create dir for tenet plugin executables. Then, work out how to install the plugins.
	defaultTenets := path.Join(home, defaultTenetCfgPath)
	if _, err := os.Stat(defaultTenets); err == nil {
		if !initPWD {
			fmt.Printf(`
error: lingo already initiated. Using these tenets:
 %s 
 If you wish to initiate this directory with it's own tenets, run:
 $ lingo init --pwd
`[1:], defaultTenets)
			return
		}

	} else {
		// init lingo

		// write tenet config
		if err := writeFileAll(defaultTenets, []byte("# run `lingo add <tenet>` to add tenets to this file"), 0777); err != nil {
			fmt.Println(err)
			return
		}

		readMe := []byte(`
This directory holds tenet executables. To add a tenet run:
$ lingo add <author>:<tenet-name>

`[1:])

		if err := writeFileAll(path.Join(home, "tenets", "README"), readMe, 0777); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("lingo successfully initiated. Tenets config file written to %q", defaultTenets)
	}

	// then init pwd
	if initPWD {
		cfgPath, err := filepath.Abs(c.GlobalString("tenet-config"))
		if err != nil {
			fmt.Printf("error: %s", err.Error())
		}

		if _, err := os.Stat(cfgPath); err == nil {
			fmt.Printf("error: pwd is already initiated. Using tenets config file %q in pwd.\n", cfgPath)
			return
		}

		if err := ioutil.WriteFile(cfgPath, []byte("# run `lingo add <tenet>` to add tenets to this file"), 0644); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("pwd successfully initiated. Tenets config file written to %q", cfgPath)
	}
}
