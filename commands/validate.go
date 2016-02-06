package commands

import (
	"github.com/lingo-reviews/lingo/commands/common"

	"github.com/codegangsta/cli"
)

var ValidateCMD = cli.Command{
	Name:  "validate",
	Usage: "validate a tenet",
	Description: `

  "lingo validate <author>/<tenet>"

`[1:],
	Action: validate,
}

// https://docs.docker.com/userguide/labels-custom-metadata/
// images should have key: reviews.lingo.tenet

func validate(c *cli.Context) {
	if err := common.ExactArgs(c, 1); err != nil {
		common.OSErrf(err.Error())
		return
	}

	// imageName := c.Args().First()
}
