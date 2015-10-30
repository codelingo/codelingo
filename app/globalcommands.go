package app

import (
	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/commands"
)

var globalCommands = []cli.Command{
	commands.InitCMD,
	commands.AddCMD,
	commands.RemoveCMD,
	commands.ImportCMD,
	commands.ReviewCMD,
	commands.PullCMD,
	commands.ListCMD,
	commands.OptionsCMD,
	commands.WhichCMD,
	commands.WriteDocCMD,
	commands.DocsCMD,
	commands.EditCMD,
	// commands.TryPieCMD,
}
