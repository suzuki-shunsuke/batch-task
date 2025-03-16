#!/usr/bin/env bash

set -eu

rm .github/workflows/wc-update-aqua-checksums.yaml
if [ -f .github/workflows/autofix.yaml ]; then
  exit 1
fi

cp "$WORKFLOW" .github/workflows/autofix.yaml
