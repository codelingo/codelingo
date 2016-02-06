package commands

import (
	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/commands/common"
)

var SearchCMD = cli.Command{
	Name:  "search",
	Usage: "search tenet image(s)",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: common.TagsFlg.String(),
		},
	},
	Description: `

lingo search <tenet-name> --tags="list,of,tags,to,search"

`[1:],
	Action: search,
}

func search(c *cli.Context) {

	// TODO(waigani) use docker client to search for tenets

}
