package config

import (
	"keys/internal/keymap"
)

type AppMode int

const (
	NormalMode AppMode = iota
	KeyboardSelectMode
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
