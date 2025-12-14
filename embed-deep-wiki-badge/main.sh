#!/usr/bin/env bash

set -eux

owner=$1
export GHTKN_APP="$owner/write"

grep -v -E "^#" "repos-${owner}.txt" | while read -r repo; do
    if [ -z "$repo" ]; then
        continue
    fi
    bash get-readme.sh "${repo}"
    if grep "img.shields.io/badge/DeepWiki-" README.md >/dev/null; then
        echo "[INFO] Skip ${repo}" >&2
        continue
    fi
    go run main.go "$repo"
    branch=docs-embed-deep-wiki-badge
    GITHUB_TOKEN=$(ghtkn get) ghcp commit \
        -r "$repo" \
        -b "$branch" \
        -m "docs: embed Deepwiki Badge" \
        README.md
    gh pr create \
        -r "$repo" \
        -B main \
        -H "$branch" \
        -t "docs: embed Deepwiki Badge" \
        -b "- https://github.com/suzuki-shunsuke/batch-task/issues/16"
done
