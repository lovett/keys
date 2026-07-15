#!/usr/bin/env sh

set -eu

cd "$(dirname "$0")/../"

. "scripts/vars.sh"

lint_js() {
	$BIOME lint internal/asset/assets/keys.js
    $BUN x tsc --noEmit
}

lint_go() {
    golangci-lint -c golangci.json run
}

lint_json() {
    jq empty < tsconfig.json
}

lint_openapi() {
    $VACUUM
}

lint_sh() {
    shellcheck scripts/*
}

case "${1-all}" in
    --help)
        echo "Run language-specific linters to check code quality."
        ;;
    go)
        lint_go
        ;;
    js)
        lint_js
        ;;
    openapi)
        lint_openapi
        ;;
    sh)
        lint_sh
        ;;
    all)
        lint_go
        lint_js
        lint_json
        lint_sh
        ;;
    *)
        echo "Unknown argument." >&2
        ;;
esac
