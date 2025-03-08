#!/usr/bin/env bash

set -eu

org=$1

export GITHUB_TOKEN=$(gh auth token)

opt=-O
if [ "$org" = suzuki-shunsuke ]; then
	opt=-U
fi

multi-gitter merge \
	-B ci-remove-create-status-check-workflow \
	--skip-forks \
	"$opt" "$org"
