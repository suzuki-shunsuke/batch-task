#!/usr/bin/env bash

set -eu

repo=$1

gh api \
    -H "Accept: application/vnd.github.object" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    "/repos/${repo}/contents/README.md" \
    --jq ".content" |
    base64 -d > README.md
