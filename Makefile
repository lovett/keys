.PHONY: build

watch:
	ls *.go assets/* | entr -r go run .

build:
	GOOS=linux GOARCH=amd64 go build

setup:
	sudo dnf install alsa-lib-devel
