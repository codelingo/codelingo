package commands

import "github.com/codegangsta/cli"

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
	exactArgs(c, 1)

	// imageName := c.Args().First()
}
