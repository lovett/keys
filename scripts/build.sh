#!/usr/bin/env sh

set -euo pipefail

cd "$(dirname "$0")/../"

if [ "$1" = "--help" ]; then
    echo "Compiles the application."
    exit
fi


go build
