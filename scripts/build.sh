#!/usr/bin/env sh

set -eu

cd "$(dirname "$0")/../"

# This value needs to work as a git tag.
VERSION=$(date +"%Y%m.%d.%H%M")

VERSION_FILE="internal/asset/assets/version.txt"

echo "$VERSION" > "$VERSION_FILE"

go build

git checkout -q "$VERSION_FILE"
