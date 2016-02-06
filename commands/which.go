package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/commands/common"
)

var WhichCMD = cli.Command{
	Name:        "which",
	Usage:       "prints path to .lingo",
	Description: "prints path to .lingo",
	Action:      which,
}

func which(c *cli.Context) {
	path, err := common.TenetCfgPath(c)
	if err != nil {
		if os.IsNotExist(err) {
			// TODO(waigani) check for error not found. Throw unexpected errors.
			fmt.Println(common.ErrMissingDotLingo.Error())
		} else {
			fmt.Println(err)
		}
		return
	}
	fmt.Println(path)
}
