#!/usr/bin/env sh

BUN_IMAGE="docker.io/oven/bun:alpine"
export BUN="podman run --rm -v $PWD:/app:Z -w /app $BUN_IMAGE bun"
BIOME_IMAGE="ghcr.io/biomejs/biome"
export BIOME="podman run --rm -v $PWD:/app:Z -w /app $BIOME_IMAGE"
VACUUM_IMAGE="docker.io/dshanley/vacuum"
export VACUUM="podman run --rm -v $PWD:/work:Z $VACUUM_IMAGE lint internal/asset/assets/openapi.yaml"
