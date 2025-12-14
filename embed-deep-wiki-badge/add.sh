#!/usr/bin/env bash

set -eu

if grep "img.shields.io/badge/DeepWiki-" README.md >/dev/null; then
    echo "[INFO] Batch has alredy been embedded" >&2
    exit 0
fi

embed-deep-wiki-badge "$REPOSITORY"
