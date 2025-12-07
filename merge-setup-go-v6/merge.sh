#!/usr/bin/env bash

set -eu

org=$1

export GITHUB_TOKEN=$(gh auth token)

opt=-O
if [ "$org" = suzuki-shunsuke ]; then
	opt=-U
fi

merge() {
	local branch=$1
	multi-gitter merge \
		-B "$branch" \
		--skip-forks \
		"$opt" "$org"
}

# merge renovate/suzuki-shunsuke-renovate-config-3.x
# merge renovate/suzuki-shunsuke-renovate-config-nolimit-3.x
merge renovate/suzuki-shunsuke-renovate-config-action-go-version-3.x
