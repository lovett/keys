#!/usr/bin/env sh

set -eu

COVERAGE_DIR="${COVERAGE_DIR-/srv/www/coverage}"

# Run from repository root
cd "$(dirname "$0")/../"

if [ "${1-}" = "--help" ]; then
    echo "Run unit tests."
    echo ""
    echo "Flags are passed on to 'go test'. See 'go help test' for details."
    exit
fi

if [ -z "$COVERAGE_DIR" ]; then
    go test -failfast -timeout 1m ./... "$@"
else
    COVERAGE_FILE="$COVERAGE_DIR/keys.html"
    go test -failfast -timeout 1m ./... -coverprofile coverage.out "$@"
    grep -v -e "testhelper.go" -e "main.go" coverage.out > coverage_filtered.out
    mv coverage_filtered.out coverage.out

    go tool cover -html=coverage.out -o="$COVERAGE_FILE"
    echo "Wrote coverage file to $COVERAGE_FILE"
fi
