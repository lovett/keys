package main

import (
	"keys/cli"
	"os"
)

const APP_VERSION = "dev"

func main() {
	os.Exit(cli.Run(APP_VERSION))
}
