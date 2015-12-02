package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/util"
)

var SetupAutoCompleteCMD = cli.Command{
	Name:  "setup-auto-complete",
	Usage: "setup auto completion lingo commands",
	Description: `
This command appends the following line to the end of ~/.bashrc:

PROG=lingo source ~/.lingo_home/scripts/bash_autocomplete.sh

That line sources an auto-complete script for lingo. To complete the setup, run:

. ~/.bashrc
lingo --generate-bash-completion
`[1:],
	Action: sourceAutoComplete,
}

func sourceAutoComplete(c *cli.Context) {
	uHome, err := util.UserHome()
	if err != nil {
		oserrf(err.Error())
		return
	}
	autoCompScriptPath := filepath.Join(uHome, ".lingo_home/scripts/bash_autocomplete.sh")
	err = ioutil.WriteFile(autoCompScriptPath, []byte(bash_autocomplete), 0775)
	if err != nil {
		fmt.Printf("WARNING: could not write script: %v \n", err)
	}

	// TODO(waigani) this just works for linux
	bashrcPath := filepath.Join(uHome, ".bashrc")
	f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		oserrf(err.Error())
		return
	}
	defer f.Close()

	if _, err = f.WriteString("PROG=lingo source " + autoCompScriptPath); err != nil {
		oserrf(err.Error())
		return
	}

	fmt.Print(`
Success! Please run the following commands to complete the setup:
. ~/.bashrc
lingo --generate-bash-completion
`[1:])
}

// --- scripts ---

var bash_autocomplete = `
#! /bin/bash

# to install, append the following line to ~/.bashrc
# PROG=lingo source ~/.lingo_home/scripts/bash_autocomplete.sh

: ${PROG:=$(basename ${BASH_SOURCE})}

_cli_bash_autocomplete() {
     local cur opts base
     COMPREPLY=()
     cur="${COMP_WORDS[COMP_CWORD]}"
     opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
     COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
     return 0
 }
  
 complete -F _cli_bash_autocomplete $PROG
 `[1:]
