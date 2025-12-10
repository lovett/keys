#!/usr/bin/env sh

set -euo pipefail

cd "$(dirname "$0")/../"

if [ "$1" = "--help" ]; then
    echo "Run language-specific linters to check code quality."
    exit
fi

lint_js() {
    biome lint internal/asset/assets/keys.js
}

lint_go() {
    golangci-lint run
}

lint_openapi() {
    vacuum dashboard --watch internal/asset/assets/openapi.yaml
}

case "$1" in
    go)
        lint_go
        ;;
    js)
        lint_js
        ;;
    openapi)
        lint_openapi
        ;;
    *)
        lint_go
        lint_js
        lint_openapi
        ;;
esac
