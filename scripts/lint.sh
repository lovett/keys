#!/usr/bin/env sh

set -eu

. "$(dirname "$0")/vars.sh"

cd "$(dirname "$0")/../"

lint_js() {
    KEYS_JS="internal/asset/assets/keys.js"
	$BIOME lint "$KEYS_JS"
    $BUN x tsc --noEmit "$KEYS_JS"
}

lint_go() {
    golangci-lint --enable=gosec run
}

lint_json() {
    if command -v "jq" > /dev/null 2>&1; then
        jq empty < tsconfig.json
    else
        echo "jq is not installed"
        exit 1
    fi
}

lint_openapi_watch() {
    vacuum dashboard --watch internal/asset/assets/openapi.yaml
}

lint_openapi() {
    case "$(vacuum lint --no-banner -q internal/asset/assets/openapi.yaml)" in
        *"100/100"*)
            echo "No issues in openapi.yaml"
            ;;
        *)
            echo "$RESULT"
            ;;
    esac
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
        lint_openapi_watch
        ;;
    sh)
        lint_sh
        ;;
    all)
        lint_go
        lint_js
        lint_json
        lint_openapi
        lint_sh
        ;;
    *)
        echo "Unknown argument." >&2
        ;;
esac
