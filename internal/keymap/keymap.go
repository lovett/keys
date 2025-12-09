package keymap

import (
	"fmt"
	"keys/internal/asset"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

var keyCache = make(map[string]*Key)

type Keymap struct {
	Filename string
	Content  *ini.File
}

func NewKeymap(filename string) (*Keymap, error) {
	km := Keymap{
		Filename: filename,
	}

	err := km.Load()
	if err != nil {
		return nil, err
	}

	return &km, nil
}

func (km *Keymap) Load() error {
	clear(keyCache)

	options := ini.LoadOptions{
		SkipUnrecognizableLines: true,
		AllowShadows:            true,
	}

	bytes, err := ini.LoadSources(options, km.Raw())
	if err != nil {
		return err
	}

	km.Content = bytes
	km.Content.BlockMode = false

	return nil
}

func (km *Keymap) Raw() []byte {
	if _, statErr := os.Stat(km.Filename); os.IsNotExist(statErr) {
		return asset.ReadKeymapSkeleton()
	}

	bytes, err := os.ReadFile(km.Filename)
	if err != nil {
		return []byte{}
	}
	return bytes
}

func (km *Keymap) FindKey(target string) *Key {
	if key, found := keyCache[target]; found {
		return key
	}

	key := km.findKeyByName(target)
	if key == nil {
		key = km.findKeyByPhysicalKey(target)
	}

	keyCache[target] = key
	return key
}

func (km *Keymap) findKeyByName(name string) *Key {
	section, err := km.Content.GetSection(name)
	if err != nil {
		return nil
	}

	return NewKeyFromSection(section, "")
}

func (km *Keymap) findKeyByPhysicalKey(physicalKey string) *Key {
	wanted := strings.ReplaceAll(physicalKey, "KEY_", "")
	wanted = strings.ReplaceAll(wanted, ",", "")

	for _, section := range km.Content.Sections() {
		iniKey, err := section.GetKey("physical_key")
		if err != nil {
			continue
		}
		if iniKey.Value() == wanted {
			return NewKeyFromSection(section, "")
		}
	}

	return nil
}

func (km *Keymap) IsPhysicalKeyPrefix(prefix string) bool {
	for _, section := range km.Content.Sections() {
		physicalKey := section.Key("physical_key").MustString("")
		if strings.HasPrefix(physicalKey, prefix) && len(prefix) < len(physicalKey) {
			return true
		}
	}
	return false
}

func (km *Keymap) Keys() func(yield func(*Key) bool) {
	row := ""

	return func(yield func(*Key) bool) {
		for _, s := range km.Content.Sections() {
			if s.Name() == ini.DefaultSection {
				continue
			}

			if strings.HasPrefix(s.Name(), "--") {
				row = strings.Trim(s.Name(), "-")
				continue
			}

			key := NewKeyFromSection(s, row)

			if key == nil {
				continue
			}

			if !yield(key) {
				return
			}
		}
	}
}

func (km *Keymap) SetKeyboard(path string) {
	km.Content.Section(ini.DefaultSection).Key("keyboard").SetValue(path)
}

func (km *Keymap) Write() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Not using system temp dir because it could on a different partition.
	// Rename across filesystems isn't supported.
	tempFile, err := os.CreateTemp(cwd, "keys-temp*.ini")
	if err != nil {
		return err
	}
	defer func() {
		removeErr := os.Remove(tempFile.Name())
		if removeErr != nil {
			err = removeErr
		}
	}()

	err = km.Content.SaveTo(tempFile.Name())
	if err != nil {
		return fmt.Errorf("could not save file %s: %w", tempFile.Name(), err)
	}

	if err := os.Rename(tempFile.Name(), km.Filename); err != nil {
		return fmt.Errorf("could not open file %s: %w", tempFile.Name(), err)
	}

	return nil
}
