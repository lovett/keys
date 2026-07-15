#!/usr/bin/env sh

set -eu

COVERAGE_DIR="${COVERAGE_DIR-/srv/www/coverage}"

# Run from repository root
cd "$(dirname "$0")/../"

run_test() {
    go test -timeout 1m ./... -coverprofile coverage.out "$@"
    grep -v -e "testhelper.go" -e "main.go" coverage.out > coverage_filtered.out
    mv coverage_filtered.out coverage.out
    go tool cover -html=coverage.out -o=$COVERAGE_DIR/keys.html
    echo "Wrote coverage file to $COVERAGE_DIR/keys.html"
}

case "${1:-default}" in
    --help)
        echo "Run unit tests."
        echo ""
        echo "Flags:"
        echo "  --verbose or -v: Pass -v to go test"
        ;;
    -v)
        run_test "-v"
        ;;
    --verbose)
        run_test "-v"
        ;;
    default)
        run_test
        ;;
    *)
        echo "Unknown argument." >&2
        ;;
esac
