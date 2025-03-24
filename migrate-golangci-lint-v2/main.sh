#!/usr/bin/env bash

set -eu

if [ ! -f .golangci.yml ]; then
    exit 0
fi

aqua -c "$WD/aqua.yaml" exec -- pinact run -u 

aqua up golangci-lint
aqua upc -prune
golangci-lint migrate
rm .golangci.bck.yml
