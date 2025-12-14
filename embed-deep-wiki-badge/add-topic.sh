#!/usr/bin/env bash

set -eu

repo=$1

export GHTKN_APP=lintnet-modules/write
gh repo edit \
    "${repo}" \
    --add-topic oss
