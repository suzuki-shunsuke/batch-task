#!/usr/bin/env bash

set -eu

git grep -l autofix-ci/action@ | xargs -n 1 sed -i -E 's|autofix-ci/action@2891949f3779a1cafafae1523058501de3d4e944 # v1.3.1|suzuki-shunsuke/autofixci-action@a148ca710fd8bf5d9959d5d7604018ce8a77bad0 # v1.3.2-1|'
