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
"copr" stands for Checkout Pull Request.

From the root of a git repository run the following:

$ lingo copr <user>:<branch>

Where <user> is the owner of the pull request and <branch> is the branch name
of the pull request. The pull request will be pulled into a new local branch
of --base with the name <user>-<branch>. Any commits back to the fork point
from base will be reset, leaving any changes in the working tree.

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
			Name:  "skip-fetch-all,s",
			Usage: "By default, all remotes are fetched before checking out the pull request, so the correct fork point can be found. Use this flag to skip the fetch.",
		},
	},
	Action: copr,
}

func copr(c *cli.Context) {
	cliArgs := c.Args()
	if err := common.ExactArgs(c, 1); err != nil {
		common.OSErrf(err.Error())
		return
	}

	repoPath, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		common.OSErrf(err.Error())
		return
	}
	repo := path.Base(strings.Trim(string(repoPath), "\n"))
	var fetchAll string
	if !c.Bool("skip-fetch-all") {
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
