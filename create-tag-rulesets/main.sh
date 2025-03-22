#!/usr/bin/env bash

set -eu

gh api \
  --method POST \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  /repos/OWNER/REPO/rulesets \
   -f "name=super cool ruleset" -f "target=branch" -f "enforcement=active" -F "bypass_actors[][actor_id]=234" -f "bypass_actors[][actor_type]=Team" -f "bypass_actors[][bypass_mode]=always" -f "ref_name[include[]]=refs/heads/main" -f "ref_name[include[]]=refs/heads/master" -f "ref_name[exclude[]]=refs/heads/dev*" -f "rules[][type]=commit_author_email_pattern" -f "rules[][parameters][operator]=contains" -f "rules[][parameters][pattern]=github"
