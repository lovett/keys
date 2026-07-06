#!/usr/bin/env sh

set -eu

cd "$(dirname "$0")/../"

VERSION=$(date +"%Y%m.%d.%H%M")
VERSION_FILE="internal/asset/assets/version.txt"

echo "$VERSION" > "$VERSION_FILE"

sh scripts/setup.sh

go build

git checkout -q "$VERSION_FILE"
