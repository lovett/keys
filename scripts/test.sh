#!/usr/bin/env sh

set -euo pipefail

cd "$(dirname "$0")/../"

if [ "$1" = "--help" ]; then
    echo "Run unit tests."
    exit
fi

go test ./internal/*
