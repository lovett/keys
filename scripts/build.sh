#!/usr/bin/env sh

set -eu

cd "$(dirname "$0")/../"

case "${1:-default}" in
    --help)
        echo "Compiles the application."
        ;;
    default)
        go build
        ;;
    *)
        echo "Unknown argument." >&2
        ;;
esac
