#!/usr/bin/env bash

set -x

path=$GOPATH
if [[ -z $path ]]; then
path="~/go"

fi

owner="$1"
name="$2"
codelingoPath="$path/src/github.com/codelingo/codelingo"
flowPath="$codelingoPath/flows/$owner/$name"

v="
windows,386;\
linux,386;\
windows,amd64;\
linux,amd64;\
darwin,amd64;"

version="0.0.0"

# Build and push each bin to release
echo $v | while IFS=',' read -d';' os arch;  do 
	
	cd $flowPath

    if [ "$os" == "windows" ]; then
    ext=.exe
    else
    ext=""
    fi
    binpath=./bin/$os/$arch/$version
   	mkdir -p $binpath
	GOOS=$os GOARCH=$arch go build -o $binpath/cmd$ext
	cd $binpath
	filename=cmd$ext
	if [ "$os" == "windows" ]; then
		fn="$filename.zip"
        rm $binpath/$fn
		zip $fn cmd.exe
        rm cmd.exe
	else
		fn="$filename.tar.gz"
		tar -cvzf $fn cmd
	    rm cmd
    fi
    cd ..
done