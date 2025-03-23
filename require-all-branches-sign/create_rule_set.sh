#!/usr/bin/env bash

set -eu

repo=$1

# gh api \
#   --method POST \
#   -H "Accept: application/vnd.github+json" \
#   -H "X-GitHub-Api-Version: 2022-11-28" \
#   "/repos/$repo/rulesets" \
#    -f "name=forbid_change" \
#    -f "target=tag" \
#    -f "enforcement=active" \
#    -f 'conditions[ref_name][include[]]=~ALL' \
#    -f "rules[][type]=deletion" \
#    -f "rules[][type]=update"

curl -L \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer $(gh auth token)" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/$repo/rulesets \
  -d '{
  "name": "require_sign",
  "target": "branch",
  "enforcement": "active",
  "conditions": {
    "ref_name": {
      "exclude": [],
      "include": [
        "~ALL"
      ]
    }
  },
  "rules": [
    {
      "type": "required_signatures"
    }
  ]
}'
