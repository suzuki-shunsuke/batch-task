#!/usr/bin/env bash

set -eu

for file in .goreleaser.yml .goreleaser.yaml; do
  if [ -f "$file" ]; then
    sed -i 's|format: zip|formats: [zip]|' "$file"
  fi
done

goreleaser check
