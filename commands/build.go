package commands

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/lingo-reviews/lingo/commands/common"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	bt "github.com/lingo-reviews/lingo/tenet/build"
	"github.com/lingo-reviews/lingo/util"
)

var allDrivers = []string{"binary", "docker"}

var BuildCMD = cli.Command{
	Name:  "build",
	Usage: "build a tenet from source",
	Description: `
	
Call "lingo build" in the root directory of the source code of a tenet.
It will look for a .lingofile with instructions on how to build your tenet.
For example:

language = "go"
owner = "lingoreviews"
name = "simpleseed"

[binary]
  build=false

[docker]
  overwrite_dockerfile=true

You can specify which driver to build. If no arguments are supplied, lingo
will try to build every driver. 
`[1:],
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "all",
			Usage: "build every tenet found in every subdirectory",
		},
		// TODO(waigani) implement
		// cli.BoolFlag{
		// 	Name:  "watch",
		// 	Usage: "watch for any changes and attempt to rebuild the tenet",
		// },
	},
	Action: buildAction,
	BashComplete: func(c *cli.Context) {
		// This will complete if no args are passed
		if len(c.Args()) > 0 {
			return
		}

		for _, d := range allDrivers {
			util.Println(d)
		}
	},
}

var c int

func buildAction(ctx *cli.Context) {
	drivers := ctx.Args()
	if len(drivers) == 0 {
		drivers = allDrivers
	}

	if err := build(drivers, ctx.Bool("all")); err != nil {
		common.OSErrf(err.Error())
	}
}

func build(drivers []string, all bool) error {
	lingofiles, err := getLingoFiles(all)
	if err != nil {
		return errors.Trace(err)
	}
	fn := len(lingofiles)
	if fn == 0 {
		return errors.Errorf("tenet not built. No %s found.\n", common.Lingofile)
	}

	if err := bt.Run(drivers, lingofiles...); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func getLingoFiles(all bool) (lingofiles []string, err error) {
	if all {
		// util.NoPrint = true

		lingofiles, err = allLingoFiles(".")
		if err != nil {
			return nil, err
		}
	} else {
		if _, err := os.Stat(common.Lingofile); err == nil {
			lingofiles = []string{common.Lingofile}
		}
	}
	return
}

func allLingoFiles(rootDir string) (lingofiles []string, err error) {
	err = filepath.Walk(rootDir, func(relPath string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(relPath, common.Lingofile) {
			lingofiles = append(lingofiles, relPath)
		}
		return nil
	})
	return
}
