#!/usr/bin/env bash

set -eu

org=$1

body="- https://github.com/suzuki-shunsuke/batch-task/issues/2"

export GITHUB_TOKEN=$(gh auth token)

opt=-O
if [ "$org" = suzuki-shunsuke ]; then
	opt=-U
fi

multi-gitter run ./replace.sh \
	--config config.yaml \
	-t "ci: use required-status-check-action" \
	-B ci-use-required-status-check-action \
	-b "$body" \
	--skip-forks \
	"$opt" "$org"
