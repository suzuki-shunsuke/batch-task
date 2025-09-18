#!/usr/bin/env bash

set -eu

perl -pi -e 's#\Qif system_command("/usr/bin/xattr", args: ["-h"]).exit_status == 0\E#if OS.mac?#g' .goreleaser.yml
