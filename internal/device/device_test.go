package device

import (
	"fmt"
	"keys/internal/config"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/holoplot/go-evdev"
)

func loadConfig(t *testing.T) *config.Config {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	fixture := filepath.Join(wd, "../../testdata/key-multiple.ini")
	cfg, err := config.NewConfig(fixture)
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func TestEcho(t *testing.T) {
	cfg := loadConfig(t)

	inputEvent := &evdev.InputEvent{
		Time:  syscall.Timeval{Sec: 0},
		Type:  evdev.EV_KEY,
		Code:  0,
		Value: 0,
	}

	deviceEvent := &DeviceEvent{"/my/device/path", inputEvent}

	want := fmt.Sprintf(`
Code: %s
From: %s
Translated to: reserved
Mapped to: none
`, inputEvent.CodeName(), deviceEvent.DevicePath)

	got := echo(deviceEvent, cfg)

	if !strings.Contains(got, want) {
		t.Errorf("wanted:\n%s\ngot:\n%s", want, got)
	}
}

func TestCanListen(t *testing.T) {
	tests := []struct {
		user   string
		group  string
		result bool
	}{
		{"nobody", "root", false},
		{"nobody", "nobody", true},
	}

	for _, tt := range tests {
		u, err := user.Lookup(tt.user)
		if err != nil {
			t.Fatal(err)
		}

		g, err := user.LookupGroup(tt.group)
		if err != nil {
			t.Fatal(err)
		}

		result := canListen(u, g)

		if result != tt.result {
			t.Errorf("Expected %t for user %s in group %s, got %t", tt.result, tt.user, tt.group, result)
		}
	}
}
