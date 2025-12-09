package config

import (
	"keys/internal/keymap"
	"os"
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
