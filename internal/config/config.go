package config

import (
	"keys/internal/keymap"
	"os"
)

type Config struct {
	KeyboardFound  bool
	KeyboardLocked bool
	Keymap         *keymap.Keymap
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
	}

	return &cfg, nil
}
