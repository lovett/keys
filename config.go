package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type AppMode int

const (
	NormalMode AppMode = iota
	VersionMode
	KeyboardSelectMode
	SystemdSetupMode
)

type Config struct {
	Mode           AppMode
	ShowAppVersion bool
	SelectKeyboard bool
	UseServer      bool
	UseKeyboard    bool
	ServerAddress  string
	PublicUrl      string
	AppVersion     string
	Keymap         *Keymap
}

func NewConfig(appVersion string) *Config {
	systemd := flag.Bool("systemd", false, "Install a systemd user service")
	version := flag.Bool("version", false, "Application version")
	selectKeyboard := flag.Bool("select-keyboard", false, "Choose which keyboard to use for input")
	inputs := flag.String("inputs", "browser,keyboard", "Where to listen for input")
	config := flag.String("config", "keys.ini", "Configuration file")
	port := flag.Int("port", 4004, "Server port")
	publicUrl := flag.String("url", "http://localhost:4004", "Server URL")

	flag.Usage = usage
	flag.Parse()

	var appMode AppMode
	if *version {
		appMode = VersionMode
	} else if *selectKeyboard {
		appMode = KeyboardSelectMode
	} else if *systemd {
		appMode = SystemdSetupMode
	} else {
		appMode = NormalMode
	}

	return &Config{
		Mode:          appMode,
		UseServer:     strings.Contains(*inputs, "browser"),
		UseKeyboard:   strings.Contains(*inputs, "keyboard"),
		ServerAddress: fmt.Sprintf(":%d", *port),
		Keymap:        NewKeymap(*config),
		AppVersion:    appVersion,
		PublicUrl:     *publicUrl,
	}
}

func (c *Config) RelativeTriggerUrl(key string) string {
	return fmt.Sprintf("/trigger/%s", key)
}

func (c *Config) PublicTriggerUrl(key string) string {
	return fmt.Sprintf("%s/trigger/%s", c.PublicUrl, key)
}

func usage() {
	fmt.Fprint(os.Stderr, "Send keyboard input to commands or services.\n\n")

	fmt.Fprint(os.Stderr, "Options\n")
	flag.PrintDefaults()
}
