#!/usr/bin/env sh

set -euo pipefail

cd "$(dirname "$0")/../"

if [ "$1" = "--help" ]; then
    echo "Install required system packages."
    exit
fi

sudo dnf install alsa-lib-devel golangci-lint
