#!/usr/bin/env sh

set -eu

cd "$(dirname "$0")/../"

sh scripts/setup.sh

go build
