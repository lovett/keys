#!/usr/bin/env sh

# shellcheck disable=SC2034
if command -v podman >/dev/null 2>&1; then
    GOLANGCI_LINT="podman run --rm -v $PWD:/app:Z -w /app docker.io/golangci/golangci-lint:latest-alpine golangci-lint"
    BIOME="podman run --rm -v $PWD:/app:Z -w /app ghcr.io/biomejs/biome"
    BUN="podman run --rm -v $PWD:/app:Z -w /app docker.io/oven/bun:1.4.2-alpine bun"
    SHELLCHECK="podman run --rm -v $PWD:/mnt:Z docker.io/koalaman/shellcheck"
    VACUUM="podman run --rm -v $PWD:/work:Z docker.io/dshanley/vacuum lint internal/asset/assets/openapi.yaml"
    # There is no container for Typescript
fi

if command -v biome >/dev/null 2>&1; then
    BIOME="biome"
fi

if command -v golangci-lint >/dev/null 2>&1; then
    GOLANGCI_LINT="golangci-lint"
fi

if command -v jq >/dev/null 2>&1; then
    JQ="jq"
fi

if command -v shellcheck >/dev/null 2>&1; then
    SHELLCHECK="shellcheck"
fi

if command -v sqlfluff >/dev/null 2>&1; then
    SQLFLUFF="sqlfluff"
fi

if command -v bun >/dev/null 2>&1; then
    TSC="bun $HOME/.bun/bin/tsc"
fi
