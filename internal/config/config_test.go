package config

import (
	"os"
	"path/filepath"
	"testing"
)

func configFromFixture(t *testing.T, filename string) (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	fixture := filepath.Join(wd, "../../testdata", filename)
	return NewConfig(fixture)
}

func TestConfigExistence(t *testing.T) {
	if _, err := configFromFixture(t, "does-not-exist"); err == nil {
		t.Error("nonexistent config file was not rejected")
	}
}

func TestConfigFileValidity(t *testing.T) {
	if _, err := configFromFixture(t, "invalid.ini"); err == nil {
		t.Error("malformed config file was not rejected")
	}
}
