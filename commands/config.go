package commands

import (
	"path/filepath"

	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/commands/common"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/util"
)

var ConfigCMD = cli.Command{
	Name:  "config",
	Usage: "open lingo's configuration file",
	Description: `

  $ lingo config

`[1:],
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "editor",
			Value: "vi",
			Usage: "editor to open the config file with",
		},
	},
	Action: editConfigAction,
}

func editConfigAction(ctx *cli.Context) {
	if err := editConfig(ctx); err != nil {
		common.OSErrf(err.Error())
	}
}

func editConfig(ctx *cli.Context) error {
	lHome, err := util.LingoHome()
	if err != nil {
		return errors.Trace(err)
	}

	filename := filepath.Join(lHome, common.ConfigFile)
	editor := ctx.String("editor")

	cmd, err := util.OpenFileCmd(editor, filename, 0)
	if err != nil {
		return errors.Trace(err)
	}

	if err = cmd.Start(); err != nil {
		return errors.Trace(err)
	}
	if err = cmd.Wait(); err != nil {
		return errors.Trace(err)
	}
	return nil
}
