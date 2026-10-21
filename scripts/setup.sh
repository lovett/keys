#!/usr/bin/env sh
# shellcheck disable=SC2086 # splitting of PACKAGES is intentional

set -eu

# Fedora
if command -v dnf >/dev/null 2>&1; then
    PACKAGES="alsa-lib-devel golangci-lint jq ShellCheck"

    if ! rpm -q $PACKAGES >/dev/null 2>&1; then
        sudo dnf --assumeyes install $PACKAGES
    fi

    exit
fi

# FreeBSD
if command -v pkg >/dev/null 2>&1; then
    PACKAGES="pkgconf golangci-lint jq hs-ShellCheck"

    if ! pkg info -e $PACKAGES; then
        doas pkg install --yes $PACKAGES
    fi

    exit
fi

echo "No supported package manager found."
exit 1
