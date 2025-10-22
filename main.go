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
		return
	case KeyTestMode:
		StartKeyTest(config)
		return
	case KeyboardSelectMode:
		keyboard, err := PromptForKeyboard()

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = config.Keymap.StoreKeyboard(*keyboard)

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf("Updated %s\n", config.Keymap.Filename)
		os.Exit(0)
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
