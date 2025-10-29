package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type AppMode int

const (
	NormalMode AppMode = iota
	VersionMode
	KeyboardSelectMode
	SystemdSetupMode
	KeyTestMode
	SoundTestMode
)

type Config struct {
	AppVersion     string
	KeyboardFound  bool
	KeyboardLocked bool
	Keymap         *Keymap
	Mode           AppMode
	PublicUrl      string
	ServerAddress  string
	UseKeyboard    bool
	UseServer      bool
}

func NewConfig(appVersion string) *Config {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not determine config dir. Giving up.")
		os.Exit(1)
	}

	soundTest := flag.Bool("soundtest", false, "Test mode to see if sound works")
	keyTest := flag.Bool("keytest", false, "Test mode to see the name of a pressed key")
	systemd := flag.Bool("systemd", false, "Install a systemd user service")
	version := flag.Bool("version", false, "Application version")
	selectKeyboard := flag.Bool("select-keyboard", false, "Choose which keyboard to use for input")
	inputs := flag.String("inputs", "browser,keyboard", "Where to listen for input")
	config := flag.String("config", filepath.Join(configDir, "keys.ini"), "Configuration file")
	port := flag.Int("port", 4004, "Server port")
	publicUrl := flag.String("url", "http://localhost:4004", "Server URL")

	flag.Usage = usage
	flag.Parse()

	appMode := NormalMode
	if *version {
		appMode = VersionMode
	} else if *selectKeyboard {
		appMode = KeyboardSelectMode
	} else if *systemd {
		appMode = SystemdSetupMode
	} else if *soundTest {
		appMode = SoundTestMode
	} else if *keyTest {
		appMode = KeyTestMode
	}

	keymap, err := NewKeymap(*config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not parse config. Giving up.")
		os.Exit(1)
	}

	return &Config{
		AppVersion:    appVersion,
		Keymap:        keymap,
		Mode:          appMode,
		PublicUrl:     *publicUrl,
		ServerAddress: fmt.Sprintf(":%d", *port),
		UseKeyboard:   strings.Contains(*inputs, "keyboard"),
		UseServer:     strings.Contains(*inputs, "browser"),
	}
}

func (c *Config) RelativeTriggerUrl(key string) string {
	return fmt.Sprintf("/trigger/%s", key)
}

func (c *Config) PublicTriggerUrl(key string) string {
	return fmt.Sprintf("%s/trigger/%s", c.PublicUrl, key)
}

func (c *Config) DesignatedKeyboard() string {
	return c.Keymap.Content.Section("").Key("keyboard").String()
}

func (c *Config) SoundAllowed() bool {
	key := "sound"
	if !c.Keymap.Content.Section("").HasKey(key) {
		return true
	}

	value, err := c.Keymap.Content.Section("").Key(key).Bool()
	if err != nil {
		return true
	}

	return value
}

func usage() {
	fmt.Fprint(os.Stderr, "Send keyboard input to commands or services.\n\n")

	fmt.Fprint(os.Stderr, "Options\n")
	flag.PrintDefaults()
}
