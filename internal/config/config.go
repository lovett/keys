package config

import (
	"keys/internal/keymap"
	"os"

	"gopkg.in/ini.v1"
)

type AppMode int

const (
	NormalMode AppMode = iota
	KeyTestMode
)

type Config struct {
	KeyboardFound  bool
	KeyboardLocked bool
	Keymap         *keymap.Keymap
	Mode           AppMode
	PublicUrl      string
}

func NewConfig(configFile string) (*Config, error) {
	_, err := os.Stat(configFile)
	if err != nil {
		return nil, err
	}

	keymap, err := keymap.NewKeymap(configFile)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Keymap: keymap,
		Mode:   NormalMode,
	}

	return &cfg, nil
}

func (c *Config) EnableKeyTestMode() {
	c.Mode = KeyTestMode
}

func (c *Config) DesignatedKeyboard() string {
	return c.defaultSectionKey("keyboard").String()
}

func (c *Config) SoundAllowed() bool {
	return c.defaultSectionKey("sound").MustBool(true)
}

func (c *Config) defaultSectionKey(key string) *ini.Key {
	return c.Keymap.Content.Section(ini.DefaultSection).Key(key)
}
