#!/usr/bin/env bash

set -eu

"$USE_TEST_WORKFLOW"

# Update aqua-checksums.json if aqua packages are removed
for dir in . aqua .aqua; do
	if [ -f "$dir/aqua-checksums.json" ] && ! git diff --quiet -- "$dir"; then
		(cd "$dir" && aqua upc -prune)
	fi
done
