#!/usr/bin/env bash

set -eu

gh api \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  /user/memberships/orgs > orgs.json

jq -r ".[].organization.login" orgs.json > orgs.txt
