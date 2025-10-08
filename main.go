package main

import (
	"fmt"
)

const APP_VERSION = "dev"

func main() {
	config := NewConfig(APP_VERSION)

	switch config.Mode {
	case VersionMode:
		fmt.Println(config.AppVersion)
		return
	case KeyboardListMode:
		for _, device := range ListDevices() {
			fmt.Println(device)
		}
		return
	case SystemdSetupMode:
		if err := InstallSystemdUserService(); err != nil {
			fmt.Println(fmt.Errorf("Error: %s", err.Error()))
		}
		return
	}

	if config.UseKeyboard {
		go StartKeyboardListener(config)
	}

	if config.UseServer {
		StartServer(config)
	}
}
