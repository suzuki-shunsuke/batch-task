#!/usr/bin/env bash

set -eu

cd "$(dirname "$0")"

go build -o bin/use-test-workflow .
export USE_TEST_WORKFLOW=$PWD/bin/use-test-workflow

ghtkn exec -e "GITHUB_TOKEN:suzuki-shunsuke/write" -- multi-gitter run "$PWD/add.sh" \
	--config config.yaml \
	"$@"
