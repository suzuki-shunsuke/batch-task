#!/usr/bin/env bash

set -eu

cd "$(dirname "$0")"

go build -o bin/replace-go-test-full-workflow .
export REPLACE_GO_TEST_FULL_WORKFLOW=$PWD/bin/replace-go-test-full-workflow

ghtkn exec -e "GITHUB_TOKEN:${GHTKN_APP:-suzuki-shunsuke/write}" -- multi-gitter run "$PWD/add.sh" \
	--config config.yaml \
	"$@"
