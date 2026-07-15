#!/usr/bin/env sh

set -eu

. "$(dirname "$0")/vars.sh"

PACKAGES="alsa-lib-devel golangci-lint ShellCheck entr jq podman"

# shellcheck disable=SC2086 # because splitting of PACKAGES is intentional.
if ! rpm -q $PACKAGES >/dev/null 2>&1; then
    echo "Installing packages..."
    sudo dnf --assumeyes install $PACKAGES
fi

if ! podman image exists "$BUN_IMAGE"; then
    echo "Pulling $BUN_IMAGE"
    podman pull -q "$BUN_IMAGE"
fi

if ! podman image exists "$BIOME_IMAGE"; then
    echo "Pulling $BIOME_IMAGE"
    podman pull -q "$BIOME_IMAGE"
fi

$BUN install
