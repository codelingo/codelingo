package commands

import "github.com/codegangsta/cli"

var RemoveCMD = cli.Command{Name: "remove",
	Usage: "remove a tenet from lingo",
	Description: `

  "lingo remove github.com/lingo-reviews/unused-args"

`[1:],
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "group",
			Value: "default",
			Usage: "group to remove tenet from"},
	},
	Action: remove,
}

func remove(c *cli.Context) {
	if l := len(c.Args()); l != 1 {
		oserrf("error: expected 1 argument, got %d", l)
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

	if !cfg.HasTenet(imageName) {
		oserrf(`error: tenet "%s" not found in %q`, imageName, c.GlobalString(tenetCfgFlg.long))
		return
	}

	if err := cfg.RemoveTenet(imageName, c.String("group")); err != nil {
		oserrf(err.Error())
		return
	}

	if err := writeConfigFile(c, cfg); err != nil {
		oserrf(err.Error())
		return
	}
}
