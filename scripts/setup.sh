#!/usr/bin/env sh

set -eu

cd "$(dirname "$0")/../"

case "${1:-default}" in
    --help)
        echo "Install required system packages."
        ;;
    default)
        sudo dnf install alsa-lib-devel golangci-lint
        ;;
    *)
        echo "Unknown argument." >&2
        ;;
esac
