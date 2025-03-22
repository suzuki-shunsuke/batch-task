#!/usr/bin/env bash

set -eu

grep -v -E "^#" repos.txt | xargs -n 1 bash create_rule_set.sh
