package commands

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"

	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/commands/common/config"
	"github.com/lingo-reviews/lingo/util"
)

var AddCMD = cli.Command{
	Name:  "add",
	Usage: "add a tenet to lingo",
	Description: `

  "lingo add github.com/lingo-reviews/unused-args"

`[1:],
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "driver",
			Value: defaultDriver(),
			Usage: "driver to use for this tenet",
		},
		cli.StringFlag{
			Name:  "registry",
			Value: "hub.docker.com",
			Usage: "the registry this tenet should be pulled from",
		},
		cli.StringFlag{
			Name:  "group",
			Value: "default",
			Usage: "group to add tenet to",
		},
		cli.StringFlag{
			Name:  "options",
			Value: "",
			Usage: "a space separated list of key=value options",
		},
	},
	Action: add,
	BashComplete: func(c *cli.Context) {
		// This will complete if no args are passed
		if len(c.Args()) > 0 {
			return
		}

		tenets, err := util.BinTenets()
		if err != nil {
			log.Printf("auto complete error %v", err)
			return
		}

		for _, t := range tenets {
			fmt.Println(t)
		}

	},
}

func defaultDriver() string {
	cfg, err := config.Defaults()
	if err != nil {
		return "binary"
	}

	driver := cfg.Tenet.Driver
	if driver != "binary" && driver != "docker" {
		common.OSErrf("invalid driver default: %q", driver)
	}
	return driver
}

func add(c *cli.Context) {
	if l := len(c.Args()); l != 1 {
		common.OSErrf("expected 1 argument, got %d", l)
		return
	}

	cfgPath, err := common.TenetCfgPath(c)
	if err != nil {
		common.OSErrf("reading config file: %s", err.Error())
		return
	}
	cfg, err := common.BuildConfig(cfgPath, common.CascadeNone)
	if err != nil {
		common.OSErrf("reading config file: %s", err.Error())
		return
	}

	imageName := c.Args().First()

	groupName := c.String("group")
	g, err := cfg.FindTenetGroup(groupName)
	if err == nil && common.HasTenet(g.Tenets, imageName) {
		common.OSErrf(`error: tenet "%s" already added`, imageName)
		return
	}

	// TODO(waigani) DEMOWARE. This will panic with wrong input. Matt didn't
	// your first PR bring in options?
	opts := map[string]interface{}{}
	if optStr := c.String("options"); optStr != "" {
		// TODO: DEMOWARE. Only set one option at a time to allow spaces in value
		//for _, part := range strings.Split(optStr, " ") {
		p := strings.Split(optStr, "=")
		opts[p[0]] = p[1]
		//}
	}

	var registry string
	driver := c.String("driver")
	if driver == "docker" {
		registry = c.String("registry")
	}

	cfg.AddTenet(common.TenetConfig{
		Name:     imageName,
		Driver:   driver,
		Registry: registry,
		Options:  opts,
	}, groupName)

	if err := common.WriteConfigFile(c, cfg); err != nil {
		common.OSErrf(err.Error())
		return
	}

	// TODO(waigani) open an interactive shell, prompt to set options.
}
