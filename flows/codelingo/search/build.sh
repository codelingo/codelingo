#!/usr/bin/env bash

dep ensure -v
go build -o ./cmd
rm -rf ./vendor