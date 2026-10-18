#!/usr/bin/env sh

set -eu

cd "$(dirname "$0")/../"

. "scripts/vars.sh"

CHANGED_FILES=$(git status --porcelain)

can_skip() {
    EXT="${1:-}"
    CMD="${2:-}"
    NAME="${3:-}"

    if [ -z "${CMD:-}" ]; then
        RED=$(tput setaf 1)
        NC=$(tput sgr0)
        echo "${RED}Skipping $NAME because it is not available${NC}"
        echo ""
        return 0
    fi

    # Checkout is pristine, or --all flag is active. Check everything.
    [ -z "${CHANGED_FILES:-}" ] && return 1

    COUNT=$(git ls-files -m -o --exclude-standard "*$EXT" | wc -l)
    if [ "$COUNT" -gt 0 ]; then
        # At least one changed file
        return 1
    fi

    # No changed files
    echo "Skipping $EXT linting";
    return 0
}

lint_js() {
    can_skip ".js" "${BIOME:-}" "biome" && return
    echo "Linting js assets with biome"
    $BIOME lint internal/asset/assets/keys.js

    echo "Linting js assets with tsc"
    can_skip ".js" "${TSC:-}" "tsc" && return
    $TSC --project .

    echo ""
}

lint_go() {
    can_skip ".go" "$GOLANGCI_LINT" "golangci-lint" && return
    echo "Linting go files with golangci-lint"
    $GOLANGCI_LINT -c golangci.json run
    echo ""
}

lint_json() {
    can_skip ".json" "$JQ" "jq"&& return
    echo "Linting tsconfig.json with jq"
    $JQ empty < tsconfig.json
    echo ""
}

lint_openapi() {
    $VACUUM
}

lint_sh() {
    can_skip ".sh" "$SHELLCHECK" "shellcheck" && return
    echo "Linting scripts with shellcheck"
    $SHELLCHECK scripts/*
    echo ""
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
    json)
        lint_json
        ;;
    openapi)
        lint_openapi
        ;;
    sh)
        lint_sh
        ;;
    all)
        CHANGED_FILES=""
        lint_go
        lint_js
        lint_json
        lint_sh
        ;;
    *)
        lint_go
        lint_js
        lint_json
        lint_openapi
        lint_sh

        echo "Unknown argument." >&2
        ;;
esac
