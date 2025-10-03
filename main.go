package main

import (
	"fmt"
	"os"
)

const APP_VERSION = "dev"

func main() {
	config := NewConfig(APP_VERSION)

	switch config.Mode {
	case VersionMode:
		fmt.Println(config.AppVersion)
		os.Exit(0)
	case KeyboardListMode:
		for _, device := range ListDevices() {
			fmt.Println(device)
		}
		os.Exit(0)
	}

	if config.UseKeyboard {
		go StartKeyboardListener(config)
	}

	if config.UseServer {
		StartServer(config)
	}
}
