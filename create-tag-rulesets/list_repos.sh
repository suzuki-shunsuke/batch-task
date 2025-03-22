#!/usr/bin/env bash

set -eu

mkdir -p owners

while read -r owner; do
    gh repo list "$owner" --no-archived --visibility public -L 1000 --json nameWithOwner > "owners/$owner.json"
done < owners.txt
