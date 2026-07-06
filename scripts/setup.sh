#!/usr/bin/env sh

set -eu

PACKAGES="alsa-lib-devel golangci-lint"

# shellcheck disable=SC2086 # because splitting of PACKAGES is intentional.
if ! rpm -q $PACKAGES >/dev/null 2>&1; then
    echo "Installing packages..."
    sudo dnf --assumeyes install $PACKAGES
fi
