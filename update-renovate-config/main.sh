#!/usr/bin/env bash

set -eu

# 'github>suzuki-shunsuke/renovate-config#2.3.0',
# 'github>suzuki-shunsuke/renovate-config:nolimit#2.3.0',

file=renovate.json
sed -i -E 's|github>suzuki-shunsuke/renovate-config(\:[^#]+)?(#[0-9]+\.[0-9]+\.[0-9]+)?|github>suzuki-shunsuke/renovate-config\1#3.2.0|' "$file"
