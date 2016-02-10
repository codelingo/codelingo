package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"

	"github.com/lingo-reviews/lingo/commands/common"
)

var PullCMD = cli.Command{
	Name:  "pull",
	Usage: "pull tenet image(s) from docker hub",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  common.AllFlg.String(),
			Usage: "pull all tenets in .lingo",
		}, cli.BoolFlag{
			Name:  common.UpdateFlg.String(),
			Usage: "look to pull a newer version",
		},
		cli.StringFlag{
			Name:  common.RegistryFlg.String(),
			Value: "hub.docker.com",
			Usage: "the registry to pull from",
		},
		cli.StringFlag{
			Name:  common.DriverFlg.String(),
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
		common.OSErrf("expected %d argument(s), got %d", expectedArgs, l)
		return
	}

	if all {
		fmt.Println("pulling all tenets found in .lingo ...")
		if err := pullAll(c); err != nil {
			common.OSErrf(err.Error())
		}
		return
	}

	// TODO(waigani) It doesn't make sense to have both registry and source.
	// Either registry for docker or source for binary.
	reg := c.String("registry")
	source := c.String("source")
	driver := c.String("driver")

	if err := pullOne(c, c.Args().First(), driver, reg, source); err != nil {
		common.OSErrf(err.Error())
	}
}

// Pull all tenets from config using assigned drivers.
func pullAll(c *cli.Context) error {
	cfgPath, err := common.TenetCfgPath(c)
	if err != nil {
		return err
	}
	cfg, err := common.BuildConfig(cfgPath, common.CascadeBoth)
	if err != nil {
		return err
	}

	ts, err := common.Tenets(c, cfg)
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

func pullOne(c *cli.Context, name, driverName, registry, source string) error {
	t, err := common.NewTenet(common.TenetConfig{
		Name:     name,
		Driver:   driverName,
		Registry: registry,
		Source:   source,
	})
	if err != nil {
		return errors.Trace(err)
	}

	return t.Pull(c.Bool("update"))
}
