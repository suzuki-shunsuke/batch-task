#!/usr/bin/env bash

set -eu

workflow="$(dirname $0)/workflow.yaml"

cp "$workflow" .github/workflows/create-status-check.yaml
