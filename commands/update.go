package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var UpdateCMD = cli.Command{
	Name:  "update",
	Usage: "update any lingo related resources",
	Description: `

  Run this command after you have installed a new version of Lingo. It will
  make required updates to $LINGO_HOME

`[1:],
	Action: update,
}

// update itself does not do anything. But calling it will trigger
// ensureLingoHome() which will add any required resources to $LINGO_HOME.
func update(c *cli.Context) {
	fmt.Println("lingo updated!")
}
