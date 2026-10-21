#!/usr/bin/env sh

set -eu

# Run from repository root
cd "$(dirname "$0")/../"

COVERAGE="${COVERAGE:-yes}"

if [ "${1-}" = "--help" ]; then
    echo "Run unit tests."
    echo ""
    echo "Arguments are passed on to 'go test'. See 'go help test' for details."
    exit
fi

if [ "$COVERAGE" = "no" ]; then
    go test -failfast -timeout 1m ./... "$@"
    exit
fi

go test -failfast -timeout 1m ./... -coverprofile coverage.out "$@"
grep -v -e "testhelper.go" -e "main.go" coverage.out > coverage_filtered.out
mv coverage_filtered.out coverage.out

go tool cover -html=coverage.out -o="coverage.html"
echo "Wrote coverage to coverage.html"
