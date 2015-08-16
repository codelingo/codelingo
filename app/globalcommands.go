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
	commands.TenetsCMD,
	commands.OptionsCMD,
	commands.WhichCMD,
	// commands.TryPieCMD,
}
