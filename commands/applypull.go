package commands

import (
	"github.com/codegangsta/cli"
	"github.com/waigani/xxx"
)

var ApplyPullCMD = cli.Command{
	Name:        "apply-pull",
	Aliases:     []string{"a"},
	Usage:       "Checkout Pull Request",
	Description: "Checkout a remote and apply a PR to that remote, unstaging all changes",
	Action:      applyPull,
}

func applyPull(c *cli.Context) {
	xxx.Print("applyPull called")
}
