#!/usr/bin/env bash

set -eu

export GITHUB_TOKEN=$(gh auth token)

multi-gitter merge \
	--config config.yaml
