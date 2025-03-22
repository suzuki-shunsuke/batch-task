#!/usr/bin/env bash

set -eu

while read -r owner; do
    gh repo list "$owner" -L 1000 --json nameWithOwner > "owners/$owner.json"
done < owners.txt
