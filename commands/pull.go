package commands

import (
	"github.com/codegangsta/cli"
	"github.com/juju/errors"
)

var PullCMD = cli.Command{
	Name:  "pull",
	Usage: "pull tenet image(s) from docker hub",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  allFlg.String(),
			Usage: "pull all tenets in .lingo",
		}, cli.BoolFlag{
			Name:  updateFlg.String(),
			Usage: "look to pull a newer version",
		},
		cli.StringFlag{
			Name:  registryFlg.String(),
			Value: "hub.docker.com",
			Usage: "the registry to pull from",
		},
		cli.StringFlag{
			Name:  driverFlg.String(),
			Value: "docker",
			Usage: "the driver used to pull and run the tenet",
		},
	},
	Description: `

  pull takes one argument, the name of the docker image or a --all flag. If
  the flag is provided, 0 arguments are expected and all tenets in .lingo
  are pulled.

`[1:],
	Action: pull,
}

func pull(c *cli.Context) {
	all := c.Bool("all")
	expectedArgs := 1
	if all {
		expectedArgs = 0
	}
	if l := len(c.Args()); l != expectedArgs {
		oserrf("expected %d argument(s), got %d", expectedArgs, l)
		return
	}

	if all {
		if err := pullAll(c); err != nil {
			oserrf(err.Error())
		}
		return
	}

	reg := c.String("registry")
	driver := c.String("driver")

	if err := pullOne(c, c.Args().First(), driver, reg); err != nil {
		oserrf(err.Error())
	}
}

// Pull all tenets from config using assigned drivers.
func pullAll(c *cli.Context) error {
	cfgPath, err := tenetCfgPath(c)
	if err != nil {
		return err
	}
	cfg, err := buildConfig(cfgPath, CascadeBoth)
	if err != nil {
		return err
	}

	ts, err := tenets(c, cfg)
	if err != nil {
		return err
	}

	for _, t := range ts {
		// TODO(waigani) don't return on err, collect errs and report at end
		err = t.Pull(c.Bool("update"))
		if err != nil {
			return err
		}
	}
	return nil
}

func pullOne(c *cli.Context, name, driverName, registry string) error {
	t, err := newTenet(c, TenetConfig{
		Name:     name,
		Driver:   driverName,
		Registry: registry,
	})
	if err != nil {
		return errors.Trace(err)
	}

	return t.Pull(c.Bool("update"))
}
