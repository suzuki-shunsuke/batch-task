#!/usr/bin/env bash

set -eu

org=$1

body="- https://github.com/suzuki-shunsuke/batch-task/issues/4"

export GITHUB_TOKEN=$(gh auth token)

opt=-O
if [ "$org" = suzuki-shunsuke ]; then
	opt=-U
fi

multi-gitter run ./main.sh \
	--config config.yaml \
	-t "ci: remove a create-status-check workflow" \
	-B ci-remove-create-status-check-workflow \
	-b "$body" \
	--skip-forks \
	"$opt" "$org"
