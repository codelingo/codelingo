#!/usr/bin/env bash

set -x

path=$GOPATH
if [[ -z $path ]]; then
path="~/go"

fi


codelingoPath="$path/src/github.com/codelingo/codelingo"
v="
windows,amd64;\
linux,amd64;\
darwin,amd64;"

# TODO(waigani) cmds should define versions
version="0.0.0"

for owner in $codelingoPath/flows/*/ ; do
	for d in $owner*/ ; do

	    cd $d

		# Build and push each bin to release
		echo $v | while IFS=',' read -d';' os arch;  do 

			cd $d

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

	done
done



echo "Done!"