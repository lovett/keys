package config

import (
	"os"
	"path/filepath"
	"testing"
)

func fixturePath(t *testing.T, filename string) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testdata := filepath.Join(wd, "../../testdata")
	return filepath.Join(testdata, filename)
}

func TestConfigExistence(t *testing.T) {
	fixture := fixturePath(t, "does-not-exist")
	_, err := NewConfig(fixture)
	if err == nil {
		t.Fatal("Nonexistent config file should have been rejected")
	}
}

func TestConfigFileValidity(t *testing.T) {
	fixture := fixturePath(t, "invalid.ini")
	_, err := NewConfig(fixture)
	if err == nil {
		t.Fatal("Malformed config file should have been rejected")
	}
}

func TestEnableKeyTestMode(t *testing.T) {
	fixture := fixturePath(t, "empty.ini")
	cfg, err := NewConfig(fixture)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Mode != NormalMode {
		t.Error("Configuration did not start in normal mode")
	}

	cfg.EnableKeyTestMode()

	if cfg.Mode != KeyTestMode {
		t.Error("Configuration did not switch to key test mode")
	}
}

func TestSoundAllowed(t *testing.T) {
	tests := []struct {
		fixture string
		want    bool
	}{
		{fixture: "sound-off.ini", want: false},
		{fixture: "sound-on.ini", want: true},
		{fixture: "sound-invalid.ini", want: true},
		{fixture: "empty.ini", want: true},
	}

	for _, tt := range tests {
		fixture := fixturePath(t, tt.fixture)
		cfg, err := NewConfig(fixture)

		if err != nil {
			t.Fatal(err)
		}

		got := cfg.SoundAllowed()
		if got != tt.want {
			t.Errorf("SoundAllowed() with %s got %#v, wanted %#v", tt.fixture, got, tt.want)
		}
	}
}

func TestDesignatedKeyboard(t *testing.T) {
	tests := []struct {
		fixture string
		want    string
	}{
		{fixture: "keyboard.ini", want: "/path/to/my/keyboard"},
		{fixture: "empty.ini", want: ""},
	}
	for _, tt := range tests {
		fixture := fixturePath(t, tt.fixture)
		cfg, err := NewConfig(fixture)

		if err != nil {
			t.Fatal(err)
		}

		got := cfg.DesignatedKeyboard()
		if got != tt.want {
			t.Errorf("DesignatedKeyboard() with %s got %#v, wanted %#v", tt.fixture, got, tt.want)
		}
	}
}
