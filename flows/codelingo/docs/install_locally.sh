#!/usr/bin/env bash

dep ensure -v
go build -o ~/.codelingo/flows/codelingo/docs/cmd
rm -rf ./vendor