#!/usr/bin/env bash

set -euo pipefail

owner=$1

grep -v -E "^#" "${owner}.txt" | while read -r repo; do
  echo "[INFO] updating $repo" >&2
  env "GHTKN_APP=${owner}/admin" gh api \
    --method PATCH \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2026-03-10" \
    "/repos/${repo}" \
   -f pull_request_creation_policy=collaborators_only
done
