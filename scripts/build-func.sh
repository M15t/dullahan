#!/usr/bin/env bash

set -e
now=$(date +'%Y-%m-%dT%T%z')
version=$(git rev-parse --short HEAD)
package="dullahan/pkg/server"

for i in /functions
do
    echo i
done

# go build -a -ldflags "-X $package.version=$version -X $package.buildTime=$now" -o server cmd/api/main.go
