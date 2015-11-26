#!/bin/bash

set -e -x

pathname=`pwd`
go build -o ~/.lingo_home/tenets/lingoreviews/$(basename $pathname)