#!/usr/bin/env bash

set -x

path=$GOPATH
if [[ -z $path ]]; then
path="~/go"

fi

codelingoPath="$path/src/github.com/codelingo/codelingo"

for flowPath in $codelingoPath/flows/*/ ; do
	for d in $flowPath*/ ; do

		owner=$(echo $d | cut -d'/' -f 10)
		name=$(echo $d |cut -d'/' -f 11)

		installedFlowPath=~/.codelingo/flows/$owner/$name
		mkdir -p $installedFlowPath

	    cd $d
	    go build -o $installedFlowPath/cmd

	done
done