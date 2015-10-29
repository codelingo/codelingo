#!/bin/bash



# rebuild_lingo.sh write-docs write-docs write-docs

remote=bitbucket
reposBaseDir=~/go/src/github.com/lingo-reviews

read -p "dev branch: " dev
# dev=${dev:-master}
read -p "lingo branch: " lingo
# lingo=${lingo:-master}
read -p "tenets branch: " tenets
# tenets=${tenets:-master}

cd "$reposBaseDir/tenets"
tenetsToBuild=`ls -d */`

set -e -x

# // TODO(waigani) for loop
# update dev
if [ -n "$dev" ]; then
	cd $reposBaseDir/dev
	s=`git status -s`
	if [ -n "$s" ]; then
	    echo "dev not empty, exiting"
	    echo $s
	    exit
	fi
	git fetch $remote
	git checkout $dev
	git pull --commit $remote $dev
fi

# update lingo
cd $reposBaseDir/lingo
if [ -n "$lingo" ]; then
	s=`git status -s`
	if [ -n "$s" ]; then
	    echo "lingo not empty, exiting"
	    echo $s
	    exit
	fi
	git fetch $remote
	git checkout $lingo
	git pull --commit $remote $lingo
fi
go install

# update tenets
cd $reposBaseDir/tenets
if [ -n "$tenets" ]; then
	s=`git status -s`
	if [ -n "$s" ]; then
	    echo "tesnet not empty, exiting"
	    echo $s
	    exit
	fi
	git fetch $remote
	git checkout $tenets
	git pull --commit $remote $tenets
fi

# update binary tenets
for tenet in $tenetsToBuild
do
  dir="$reposBaseDir/tenets/"$tenet
  cd $dir
  install_tenet.sh
done