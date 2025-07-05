#!/usr/bin/env bash

set -eu

git grep -l -E " +if: failure()" .github/workflows | xargs -n 1 sed -i -E "s/failure\(\)/always() \&\& (contains(needs.*.result, 'failure') || contains(needs.*.result, 'cancelled'))/"
actionlint
