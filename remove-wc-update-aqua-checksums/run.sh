#!/usr/bin/env bash

set -eu

org=$1

export GITHUB_TOKEN=$(gh auth token)

opt=-O
if [ "$org" = suzuki-shunsuke ]; then
	opt=-U
fi

export WORKFLOW=$PWD/workflow.yaml

multi-gitter run ./add.sh \
	--config config.yaml \
	"$opt" "$org"
