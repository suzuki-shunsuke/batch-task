#!/usr/bin/env bash

set -eu

export GITHUB_TOKEN=$(gh auth token)

export WD=$(pwd)

multi-gitter run ./main.sh \
    --config config.yaml
