package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Keymap struct {
	Filename string
	Content  *ini.File
	KeyCache map[string]*Key
}

func NewKeymap(filename string) *Keymap {
	k := Keymap{
		Filename: filename,
	}

	k.Parse()
	k.KeyCache = make(map[string]*Key)
	return &k
}

func (k *Keymap) Reload() {
	k.Parse()
	k.KeyCache = make(map[string]*Key)
}

func (k *Keymap) Parse() error {
	var err error

	options := ini.LoadOptions{
		SkipUnrecognizableLines: true,
		AllowShadows:            true,
	}

	if _, statErr := os.Stat(k.Filename); os.IsNotExist(statErr) {
		skeleton, err := ReadAsset("assets/skeleton.ini")
		if err != nil {
			return err
		}

		k.Content, err = ini.LoadSources(options, skeleton.Bytes)
	} else {
		k.Content, err = ini.LoadSources(options, k.Filename)
	}

	if err != nil {
		return err
	}

	// Usage is read-only, so (maybe?) speed up read operations.
	//
	// See https://ini.unknwon.io/docs/faqs
	k.Content.BlockMode = false

	return nil
}

func (k *Keymap) Raw() []byte {
	blank := []byte{}

	if _, statErr := os.Stat(k.Filename); os.IsNotExist(statErr) {
		asset, err := ReadAsset("assets/skeleton.ini")
		if err != nil {
			return blank
		}
		return asset.Bytes
	}

	bytes, err := os.ReadFile(k.Filename)
	if err != nil {
		return blank
	}
	return bytes
}

func (k *Keymap) KeyNameToSectionName(keyName string) string {
	sectionName := strings.ReplaceAll(keyName, "KEY_", "")
	sectionName = strings.ReplaceAll(sectionName, ",", "")
	return strings.ToLower(sectionName)
}

func (k *Keymap) IsMappedKey(keyName string) bool {
	sectionName := k.KeyNameToSectionName(keyName)
	for _, value := range k.Content.SectionStrings() {
		if value == sectionName {
			return true
		}
	}
	return false
}

func (k *Keymap) NewKey(keyName string) *Key {
	if key, ok := k.KeyCache[keyName]; ok {
		return key
	}

	sectionName := k.KeyNameToSectionName(keyName)
	section, err := k.Content.GetSection(sectionName)

	if err != nil {
		return nil
	}

	key := NewKeyFromSection(section)
	k.KeyCache[keyName] = key
	return key

}

func (k *Keymap) IsPrefix(keyBuffer []string) bool {
	counter := 0
	keys := k.KeyNameToSectionName(strings.Join(keyBuffer, ","))
	for _, section := range k.Content.SectionStrings() {
		if strings.HasPrefix(section, keys) {
			counter++
		}
	}

	return counter > 1
}

func (k *Keymap) Keys() func(yield func(*Key) bool) {
	return func(yield func(*Key) bool) {
		for _, s := range k.Content.Sections() {
			if s.Name() == ini.DefaultSection {
				continue
			}

			key := k.NewKey(s.Name())

			if key == nil {
				continue
			}

			if !yield(key) {
				return
			}
		}
	}
}

func (k *Keymap) StoreKeyboard(path string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.New("Failed to get current directory.")
	}

	// Not using system temp dir because rename across filesystems isn't supported
	// and /tmp is probably on a separate partition.
	tempFile, err := os.CreateTemp(cwd, "keys-temp*.ini")
	if err != nil {
		return errors.New("Failed to create temporary file.")
	}
	defer os.Remove(tempFile.Name())

	k.Content.Section("").Key("keyboard").SetValue(path)

	k.Content.SaveTo(tempFile.Name())

	if err := os.Rename(tempFile.Name(), k.Filename); err != nil {
		return fmt.Errorf("could not open file %q: %w", tempFile.Name(), err)
	}

	return nil
}
