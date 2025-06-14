#!/usr/bin/env bash

set -eu

export GITHUB_TOKEN=$(gh auth token)
export WORKFLOW=$PWD/actionlint.yaml

multi-gitter run ./add.sh \
	--config config.yaml
