package commands

import (
	"strings"

	"github.com/codegangsta/cli"
)

var AddCMD = cli.Command{
	Name:  "add",
	Usage: "add a tenet to lingo",
	Description: `

  "lingo remove github.com/lingo-reviews/unused-args"

`[1:],
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "driver",
			Value: "docker",
			Usage: "driver to use for this tenet",
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
}

func add(c *cli.Context) {
	if l := len(c.Args()); l != 1 {
		oserrf("expected 1 argument, got %d", l)
		return
	}

	cfgPath, err := tenetCfgPath(c)
	if err != nil {
		oserrf("reading config file: %s", err.Error())
		return
	}
	cfg, err := buildConfig(cfgPath, CascadeNone)
	if err != nil {
		oserrf("reading config file: %s", err.Error())
		return
	}

	imageName := c.Args().First()

	groupName := c.String("group")
	g, err := cfg.FindTenetGroup(groupName)
	if err == nil && hasTenet(g.Tenets, imageName) {
		oserrf(`error: tenet "%s" already added`, imageName)
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

	cfg.AddTenet(TenetConfig{
		Name:    imageName,
		Driver:  c.String("driver"),
		Options: opts,
	}, groupName)

	if err := writeConfigFile(c, cfg); err != nil {
		oserrf(err.Error())
		return
	}

	// TODO(waigani) open an interactive shell, prompt to set options.
}
