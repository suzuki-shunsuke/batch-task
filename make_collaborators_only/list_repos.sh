#!/usr/bin/env bash

set -euo pipefail

for owner in lintnet lintnet-modules csm-actions aquaproj suzuki-shunsuke; do
    if [ -f "${owner}.txt" ]; then
        continue
    fi
    gh search repos topic:oss \
        --owner "$owner" \
        --archived=false \
        -L 1000 \
        --json fullName \
        --jq ".[].fullName" > "${owner}.txt"
done
