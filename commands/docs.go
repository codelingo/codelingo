package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
	"github.com/pkg/browser"
)

var DocsCMD = cli.Command{
	Name:   "docs",
	Usage:  "output documentation generated from tenets",
	Action: docs,
}

func docs(c *cli.Context) {

	fmt.Println("Opening tenet documentation in a new browser window ...")
	// TODO(waigani) DEMOWARE
	writeTenetDoc(c, "", "/tmp/tenets.md")
	err := browser.OpenFile("/tmp/tenets.md")
	if err != nil {
		log.Printf("ERROR opening docs in browser: %v", err)
	}
}
