#!/usr/bin/env sh

set -eu

cd "$(dirname "$0")/../"

case "${1:-default}" in
    --help)
        echo "Run unit tests."
        ;;
    default)
        go test ./internal/*
        ;;
    *)
        echo "Unknown argument." >&2
        ;;
esac
