package commands

import (
	"os/exec"

	"github.com/codegangsta/cli"
)

var DocsCMD = cli.Command{
	Name:   "docs",
	Usage:  "output documentation generated from tenets",
	Action: docs,
}

func docs(c *cli.Context) {

	// TODO(waigani) DEMOWARE
	writeTenetDoc(c, "", "/tmp/tenets.md")
	cmd := exec.Command("chromium-browser", "/tmp/tenets.md")
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	cmd.Run()
}
