#! /bin/bash

# to install, append the following line to ~/.bashrc
# PROG=lingo source ~/go/src/github.com/lingo-reviews/lingo/scripts/bash_autocomplete.sh

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