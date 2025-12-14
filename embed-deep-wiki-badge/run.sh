#!/usr/bin/env bash

set -eu

owner=$1

export GITHUB_TOKEN=$(ghtkn get "$owner/write")

opt=-O
if [ "$owner" = suzuki-shunsuke ]; then
	opt=-U
fi

multi-gitter run ./add.sh \
	--config config.yaml \
	"$opt" "$owner"
