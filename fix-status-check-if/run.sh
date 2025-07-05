#!/usr/bin/env bash

set -eu

export GITHUB_TOKEN=$(gh auth token)

multi-gitter run ./main.sh \
    --config config.yaml
