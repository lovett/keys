package keymap

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"gopkg.in/ini.v1"
)

func keymapFromFixture(t *testing.T, filename string) *Keymap {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testdata := filepath.Join(wd, "../../testdata")
	fixture := filepath.Join(testdata, filename)
	km, err := NewKeymap(fixture)
	if err != nil {
		t.Fatal(err)
	}

	return km
}

func clearCache() {
	clear(keyCache)
}

func TestFindKeyByName(t *testing.T) {
	t.Cleanup(clearCache)

	tests := []struct {
		needle   string
		strategy string
		match    bool
	}{
		{needle: "test", strategy: "name", match: true},
		{needle: "abc123", strategy: "name", match: false},
		{needle: "hi", strategy: "physical key", match: true},
		{needle: "x", strategy: "physical key", match: false},
	}

	km := keymapFromFixture(t, "key-single.ini")

	for _, tt := range tests {
		clear(keyCache)

		key := km.FindKey(tt.needle)

		if key == nil && tt.match {
			t.Fatalf("False negative match by %s", tt.strategy)
		}

		if key != nil && !tt.match {
			t.Fatalf("False positive match by %s", tt.strategy)
		}

		if len(keyCache) == 0 {
			t.Fatal("Key cache was empty after successful lookup")
		}
	}
}

func TestPrefixDetection(t *testing.T) {
	t.Cleanup(clearCache)

	km := keymapFromFixture(t, "key-multiple.ini")

	tests := []struct {
		needle string
		match  bool
	}{
		{needle: "h", match: true},
		{needle: "x", match: false},
		{needle: "hi", match: false}, // an exact match is not a prefix
	}

	for _, tt := range tests {
		clear(keyCache)
		result := km.IsPhysicalKeyPrefix(tt.needle)
		if result != tt.match {
			if tt.match == true {
				t.Fatalf("Prefix detection false negative for \"%s\"", tt.needle)
			} else {
				t.Fatalf("Prefix detection false positive for \"%s\"", tt.needle)
			}
		}
	}
}

func TestIteration(t *testing.T) {
	t.Cleanup(clearCache)

	km := keymapFromFixture(t, "key-multiple.ini")

	keys := slices.Collect(km.Keys())

	if len(keys) != 3 {
		t.Fatalf("Iteration length mismatch. Wanted 3, got %d", len(keys))
	}

	if keys[0].Row != "" {
		t.Fatalf("First key should not have a row. Got %s", keys[0].Row)
	}

	if keys[1].Row == "" {
		t.Fatalf("Second key should have a row.")
	}
}

func TestSetKeyboard(t *testing.T) {
	t.Cleanup(clearCache)
	km := keymapFromFixture(t, "key-multiple.ini")

	path := "/path/to/keyboard"
	km.SetKeyboard(path)

	if km.Content.Section(ini.DefaultSection).Key("keyboard").String() != path {
		t.Fatal("Keyboard path not found in default section after being set")
	}
}

func TestSave(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tempFile, err := os.CreateTemp(cwd, "keys-test-temp*.ini")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := os.Remove(tempFile.Name())
		if err != nil {
			t.Fatal(err)
		}
	})

	km := keymapFromFixture(t, "key-multiple.ini")

	needle := "/keyboard-here"

	km.SetKeyboard(needle)
	km.Filename = tempFile.Name()

	err = km.Write()
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), needle) {
		t.Fatal("Keymap file did not contain expected value after write.")
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
		km := keymapFromFixture(t, tt.fixture)

		if km.SoundAllowed != tt.want {
			t.Errorf("SoundAllowed with %s got %#v, wanted %#v", tt.fixture, km.SoundAllowed, tt.want)
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
		km := keymapFromFixture(t, tt.fixture)

		if km.DesignatedKeyboard != tt.want {
			t.Errorf("DesignatedKeyboard() with %s got %#v, wanted %#v", tt.fixture, km.DesignatedKeyboard, tt.want)
		}
	}
}

func TestTranslate(t *testing.T) {
	tests := []struct {
		before string
		after  string
	}{
		{"KEY_A", "a"},
		{"KEY_A,KEY_B", "ab"},
	}

	for _, tt := range tests {
		result := Translate(tt.before)
		if result != tt.after {
			t.Errorf("wanted %s, got %s", tt.after, result)
		}
	}
}
