#!/usr/bin/env bash

set -eu

git grep -l "runs-on: ubuntu-latest" .github/workflows |
    xargs -n 1 gsed -i "s|runs-on: ubuntu-latest|runs-on: ubuntu-26.04|g"

git grep -l "runs-on: ubuntu-24.04" .github/workflows |
    xargs -n 1 gsed -i "s|runs-on: ubuntu-24.04|runs-on: ubuntu-26.04|g"
