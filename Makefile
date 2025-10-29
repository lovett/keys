.PHONY: build lint lint-go lint-js setup watch

build:
	go build

lint: lint-go lint-js

lint-go:
	golangci-lint run

lint-js:
	biome lint assets/keys.js || true

setup:
	sudo dnf install alsa-lib-devel golangci-lint

watch:
	ls *.go assets/* | entr -r go run .
