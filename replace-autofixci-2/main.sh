#!/usr/bin/env bash

set -eu

git grep -l suzuki-shunsuke/autofixci-action@ | xargs -n 1 sed -i -E 's|suzuki-shunsuke/autofixci-action@.* # .*$|autofix-ci/action@635ffb0c9798bd160680f18fd73371e355b85f27 # v1.3.2|'
