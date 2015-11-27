package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var ApplyPullCMD = cli.Command{
	Name:        "apply-pull",
	Aliases:     []string{"a"},
	Usage:       "Checkout Pull Request",
	Description: "Checkout a remote and apply a PR to that remote, unstaging all changes",
	Action:      applyPull,
}

func applyPull(c *cli.Context) {
	fmt.Println("apply-pull not implemented")
}
