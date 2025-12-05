.PHONY: build lint lint-go lint-js lint-openapi setup watch

run:
	go run . start

build:
	go build

lint: lint-go lint-js lint-openapi

lint-go:
	golangci-lint run

lint-js:
	biome lint internal/asset/assets/keys.js

lint-openapi:
	vacuum dashboard --watch internal/asset/assets/openapi.yaml

setup:
	sudo dnf install alsa-lib-devel golangci-lint

watch:
	find ./internal -type f | entr -r make run
