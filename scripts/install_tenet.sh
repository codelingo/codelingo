#!/bin/bash

set -e -x

pathname=`pwd`
go build -o ~/.lingo/tenets/lingoreviews/$(basename $pathname)