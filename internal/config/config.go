package config

import (
	"flag"
	"fmt"
	"keys/internal/keymap"
	"os"
)

type AppMode int

const (
	NormalMode AppMode = iota
	VersionMode
	KeyboardSelectMode
	KeyTestMode
)

type Config struct {
	AppVersion     string
	KeyboardFound  bool
	KeyboardLocked bool
	Keymap         *keymap.Keymap
	Mode           AppMode
	PublicUrl      string
}

func NewConfig(configFile string, appVersion string) *Config {
	selectKeyboard := flag.Bool("select-keyboard", false, "Choose which keyboard to use for input")

	flag.Parse()

	appMode := NormalMode
	if *selectKeyboard {
		appMode = KeyboardSelectMode
	}

	keymap, err := keymap.NewKeymap(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not parse config. Giving up.")
		os.Exit(1)
	}

	return &Config{
		AppVersion: appVersion,
		Keymap:     keymap,
		Mode:       appMode,
	}
}

func (c *Config) EnableKeyTestMode() {
	c.Mode = KeyTestMode
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
