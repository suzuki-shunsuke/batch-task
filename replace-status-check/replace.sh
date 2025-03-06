#!/usr/bin/env bash

set -euo pipefail

git ls-files | grep -E "^\.github/workflows/" | xargs -n 1 -r sed -E -i ':a;N;$!ba;s/    if: failure\(\)\n(.*\n)?    steps:\n      - run: exit 1/    if: always()\n\1    steps:\n      - uses: suzuki-shunsuke\/required-status-check-action@2b5a46064846b09381852c2c4217e898f639e768 # v0.1.3\n        with:\n          needs: ${{ toJson(needs) }}/'
