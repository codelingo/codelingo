package commands

import (
	"fmt"
	"os/exec"

	"github.com/codegangsta/cli"
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
	cmd := exec.Command("chromium-browser", "/tmp/tenets.md")
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	cmd.Run()
}
