package app

import (
	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/commands"
)

var globalCommands = []cli.Command{
	commands.InitCMD,
	commands.AddCMD,
	commands.RemoveCMD,
	commands.ReviewCMD,
	commands.PullCMD,
	commands.ListCMD,
	commands.WhichCMD,
	commands.WriteDocCMD,
	commands.DocsCMD,
	commands.EditCMD,
	commands.InfoCMD,
	commands.BuildCMD,
	commands.PushCMD,
	commands.SetupAutoCompleteCMD,
	commands.CoprCMD,
	commands.UpdateCMD,

	// commands.OptionsCMD,
	// commands.ImportCMD,
}
