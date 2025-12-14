#!/usr/bin/env bash

set -eu

owner=$1

gh search repos \
    --topic oss \
    --archived=false \
    --owner "${owner}" \
    --json fullName \
    --jq '.[].fullName' > "repos-${owner}.txt"
