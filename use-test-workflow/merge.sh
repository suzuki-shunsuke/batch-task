#!/usr/bin/env bash

set -eu

ghtkn exec -e "GITHUB_TOKEN:suzuki-shunsuke/write" -- multi-gitter merge \
	--config config.yaml
