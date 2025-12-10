#!/usr/bin/env sh

set -euo pipefail

cd "$(dirname "$0")/../"

if [ "$1" = "--help" ]; then
    echo "Run the application and auto restart when files change."
    exit
fi

find ./internal -type f | entr -r go run . start
