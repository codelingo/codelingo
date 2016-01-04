package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var WhichCMD = cli.Command{
	Name:        "which",
	Usage:       "prints path to .lingo",
	Description: "prints path to .lingo",
	Action:      which,
}

func which(c *cli.Context) {
	path, err := tenetCfgPath(c)
	if err != nil {
		if os.IsNotExist(err) {
			// TODO(waigani) check for error not found. Throw unexpected errors.
			fmt.Println(errMissingDotLingo.Error())
		} else {
			fmt.Println(err)
		}
		return
	}
	fmt.Println(path)
}
