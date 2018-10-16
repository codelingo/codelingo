#!/usr/bin/env bash

dep ensure -v
go build -o ~/.codelingo/flows/codelingo/search/cmd
rm -rf ./vendor