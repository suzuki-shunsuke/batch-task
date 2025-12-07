#!/usr/bin/env bash

set -eu

branch=$1

for org in lintnet lintnet-modules csm-actions aquaproj; do
    GITHUB_TOKEN=$(ghtkn get "$org/write") multi-gitter merge --org "$org" -B "$branch"
done

GITHUB_TOKEN=$(ghtkn get "suzuki-shunsuke/write") multi-gitter merge --user suzuki-shunsuke -B "$branch"
