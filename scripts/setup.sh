#!/usr/bin/env sh

set -eu

# Fedora
if command -v dnf >/dev/null 2>&1; then
    sudo dnf --assumeyes install \
         alsa-lib-devel golangci-lint jq ShellCheck
    exit
fi

# FreeBSD
if command -v pkg >/dev/null 2>&1; then
    # No package available for biome
    doas pkg install --yes \
         pkgconf golangci-lint jq hs-ShellCheck py312-sqlfluff
    exit
fi

echo "No supported package manager found."
exit 1
