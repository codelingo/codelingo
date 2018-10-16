#!/usr/bin/env bash

set -x

path=$GOPATH
if [[ -z $path ]]; then
path="~/go"

fi


codelingoPath="$path/src/github.com/codelingo/codelingo"

for owner in $codelingoPath/flows/*/ ; do
	for d in $owner*/ ; do
	    cd $d
	    go build -o ./cmd
	done
done
