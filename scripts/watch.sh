#!/usr/bin/env sh

set -eu

cd "$(dirname "$0")/../"

case "${1:-default}" in
    --help)
        echo "Run the application and auto restart when files change."
        ;;
    default)
        killall -q keys || true
        find ./internal -type f | entr -r go run . start
        ;;
    *)
        echo "Unknown argument." >&2
        ;;
esac
