#!/usr/bin/env bash

set -eu

org=$1

body="https://github.com/suzuki-shunsuke/batch-task/issues/1"

export GITHUB_TOKEN=$(gh auth token)

opt=-O
if [ "$org" = suzuki-shunsuke ]; then
	opt=-U
fi

multi-gitter run ./add.sh \
	--config config.yaml \
	-t "ci: add a workflow to check if commits in pull requests are verified" \
	-B ci-add-workflow-check-commit-verified \
	--skip-forks \
	"$opt" "$org"
