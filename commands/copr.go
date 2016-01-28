package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/lingo-reviews/lingo/commands/common"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/util"
)

var CoprCMD = cli.Command{
	Name:  "copr",
	Usage: "Checkout a git pull request",
	Description: `

cd into a local git repository the pull request targets. The run:

$ lingo copr <user>:<branch>

This will: create a new branch; checkout the pull request; and reset any
commits back to the point the branch forked from base.

`[1:],
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "base,b",
			Value: "master",
			Usage: "The base branch the pull request is forked from.",
		},
		cli.BoolFlag{
			Name:  "dry-run,d",
			Usage: "Prints out what this command will do, without doing it.",
		},
		cli.StringFlag{
			Name:  "host",
			Value: "https://github.com",
			Usage: "The remote service hosting the git repository",
		},
		cli.BoolFlag{
			Name:  "fetch-all,f",
			Usage: "Fetch all remotes before checking out pull request. This will ensure the correct fork point can be found.",
		},
	},
	Action: copr,
}

func copr(c *cli.Context) {
	cliArgs := c.Args()
	if l := len(cliArgs); l != 1 {
		common.OSErrf("expected one arguments, got %d", l)
		return
	}

	repoPath, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		common.OSErrf(err.Error())
		return
	}
	repo := path.Base(strings.Trim(string(repoPath), "\n"))
	var fetchAll string
	if c.Bool("fetch-all") {
		fetchAll = "git fetch --all"
	}

	argParts := strings.Split(cliArgs[0], ":")
	args := map[string]string{
		"User":     argParts[0],
		"Repo":     repo,
		"Branch":   argParts[1],
		"Base":     strings.Replace(c.String("base"), ":", "/", -1),
		"FetchAll": fetchAll,
		"Host":     c.String("host"),
	}
	bashScript, err := util.FormatOutput(args, coprcmd)
	if err != nil {
		common.OSErrf(err.Error())
		return
	}

	if c.Bool("dry-run") {
		fmt.Println("execute the following bash script:")
		fmt.Println(bashScript)
		return
	}

	cmd := exec.Command("bash", "-c", bashScript)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		common.OSErrf(err.Error())
		return
	}
}

var coprcmd = `
#!/bin/bash

set -e

status=` + "`git status -s`" + `
echo $status
if [ -n "$status" ]; then
	echo "aborting: working directory not clean"
	exit
fi


set -x
{{.FetchAll}}
git checkout -b {{.User}}-{{.Branch}} {{.Base}}
git pull {{.Host}}/{{.User}}/{{.Repo}}.git {{.Branch}}

sha=` + "`git merge-base --fork-point HEAD {{.Base}}`" + `

git reset $sha
`[1:]
