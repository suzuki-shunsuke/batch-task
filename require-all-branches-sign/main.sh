#!/usr/bin/env bash

set -eu

export GITHUB_TOKEN=$(gh auth token)

grep -v -E "^#" repos.txt | xargs -n 1 bash create_rule_set.sh
