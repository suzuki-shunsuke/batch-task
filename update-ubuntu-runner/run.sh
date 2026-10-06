#!/usr/bin/env bash

set -eu

owner=$1

opt=-O
if [ "$owner" = suzuki-shunsuke ]; then
	opt=-U
fi

ghtkn exec -e "GITHUB_TOKEN:$owner/write" -- multi-gitter run ./add.sh \
	--config config.yaml \
	"$opt" "$owner"
