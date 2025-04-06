#!/usr/bin/env bash

set -eu

while read -r owner; do
    jq -r ".[].nameWithOwner" "owners/$owner.json" >> repos.txt
done < owners.txt
