#!/usr/bin/env bash

set -eu

git grep -l suzuki-shunsuke/go-autofix-action@ | xargs -n 1 sed -i -E 's|suzuki-shunsuke/go-autofix-action@559f0cb21668a975222826fa376dbec951abc2f6 # v0.1.4|suzuki-shunsuke/go-autofix-action@13021ffe2dd68264a74ff8b11d3a16e19ddcc7a5 # v0.1.5|'
