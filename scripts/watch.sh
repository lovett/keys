#!/usr/bin/env sh

# Run the application and auto restart when files change.

set -eu

# Run from repository root
cd "$(dirname "$0")/../"

killall -q keys || true

export ENTR_RESTART_SIGNAL=SIGINT
find ./internal -type f | entr -r go run . start
