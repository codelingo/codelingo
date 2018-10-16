#!/usr/bin/env bash

dep ensure -v
go build -o ~/.codelingo/flows/codelingo/rewrite/cmd
rm -rf ./vendor