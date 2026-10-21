#!/usr/bin/env sh

# shellcheck disable=SC2034
if command -v podman >/dev/null 2>&1; then
    BIOME="podman run --rm -v $PWD:/app:Z -w /app ghcr.io/biomejs/biome"
    BUN="podman run --rm -v $PWD:/app:Z -w /app docker.io/oven/bun:alpine bun"
    GOLANGCI_LINT="podman run --rm -v $PWD:/app:Z -w /app docker.io/golangci/golangci-lint:latest-alpine golangci-lint"
    SHELLCHECK="podman run --rm -v $PWD:/mnt:Z docker.io/koalaman/shellcheck"
    VACUUM="podman run --rm -v $PWD:/work:Z docker.io/dshanley/vacuum lint internal/asset/assets/openapi.yaml"
fi

if command -v bun >/dev/null 2>&1; then
    BUN="bun"
fi

if command -v biome >/dev/null 2>&1; then
    BIOME="biome"
fi

if command -v golangci-lint >/dev/null 2>&1; then
    GOLANGCI_LINT="golangci-lint"
fi

if command -v shellcheck >/dev/null 2>&1; then
    SHELLCHECK="shellcheck"
fi

if command -v sqlfluff >/dev/null 2>&1; then
    SQLFLUFF="sqlfluff"
fi
