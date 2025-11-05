package main

import (
	"fmt"
	"keys/internal/config"
	"keys/internal/keyboard"
	"keys/internal/server"
	"keys/internal/sound"
	"keys/internal/system"
	"os"
)

const APP_VERSION = "dev"

func main() {
	cfg := config.NewConfig(APP_VERSION)

	switch cfg.Mode {
	case config.VersionMode:
		fmt.Println(cfg.AppVersion)
		return
	case config.KeyTestMode:
		keyboard.StartKeyTest(cfg)
		return
	case config.SoundTestMode:
		sound.StartSoundTest()
		return
	case config.KeyboardSelectMode:
		keyboard, err := keyboard.PromptForKeyboard()

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = cfg.Keymap.StoreKeyboard(*keyboard)

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf("Updated %s\n", cfg.Keymap.Filename)
		os.Exit(0)
	case config.SystemdSetupMode:
		if err := system.InstallSystemdUserService(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return
	}

	if cfg.UseKeyboard {
		go keyboard.StartKeyboardListener(cfg)
	}

	if cfg.UseServer {
		sound.LoadSounds()
		server.StartServer(cfg)
	}
}
