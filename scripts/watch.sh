#!/usr/bin/env sh

# Run the application and auto restart when files change.

set -eu

cd "$(dirname "$0")/../"

killall -q keys || true

find ./internal -type f | entr -r go run . start
