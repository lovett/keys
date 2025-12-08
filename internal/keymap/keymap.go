package keymap

import (
	"errors"
	"fmt"
	"keys/internal/asset"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Keymap struct {
	Filename string
	Content  *ini.File
	KeyCache map[string]*Key
}

func NewKeymap(filename string) (*Keymap, error) {
	k := Keymap{
		Filename: filename,
	}

	err := k.Parse()
	if err != nil {
		return nil, err
	}

	k.KeyCache = make(map[string]*Key)
	return &k, nil
}

func (k *Keymap) Reload() error {
	err := k.Parse()
	if err != nil {
		return err
	}

	k.KeyCache = make(map[string]*Key)
	return nil
}

func (k *Keymap) Parse() error {
	var loadErr error

	options := ini.LoadOptions{
		SkipUnrecognizableLines: true,
		AllowShadows:            true,
	}

	if _, statErr := os.Stat(k.Filename); os.IsNotExist(statErr) {
		skeleton, err := asset.ReadAsset("assets/skeleton.ini")
		if err != nil {
			return err
		}

		k.Content, loadErr = ini.LoadSources(options, skeleton.Bytes)
	} else {
		k.Content, loadErr = ini.LoadSources(options, k.Filename)
	}

	if loadErr != nil {
		return loadErr
	}

	// Usage is read-only, so (maybe?) speed up read operations.
	// See https://ini.unknwon.io/docs/faqs
	k.Content.BlockMode = false

	return nil
}

func (k *Keymap) Raw() []byte {
	blank := []byte{}

	if _, statErr := os.Stat(k.Filename); os.IsNotExist(statErr) {
		asset, err := asset.ReadAsset("assets/skeleton.ini")
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

func (k *Keymap) KeyNameToPhysicalKey(keyName string) string {
	sectionName := strings.ReplaceAll(keyName, "KEY_", "")
	sectionName = strings.ReplaceAll(sectionName, ",", "")
	return strings.ToLower(sectionName)
}

func (k *Keymap) IsMappedKey(name string) bool {
	physicalKey := k.KeyNameToPhysicalKey(name)
	return k.NewKey(physicalKey) == nil
}

func (k *Keymap) NewKey(name string) *Key {
	if key, ok := k.KeyCache[name]; ok {
		return key
	}

	var key *Key

	section, err := k.Content.GetSection(name)
	if err == nil {
		key = NewKeyFromSection(section)
	} else {
		physicalKey := k.KeyNameToPhysicalKey(name)

		for _, section = range k.Content.Sections() {
			iniKey, err := section.GetKey("physical_key")
			if err != nil {
				continue
			}
			if iniKey.Value() == physicalKey {
				key = NewKeyFromSection(section)
				break
			}
		}
	}

	k.KeyCache[name] = key
	return key
}

func (k *Keymap) IsPrefix(keyBuffer []string) bool {
	counter := 0
	keys := k.KeyNameToPhysicalKey(strings.Join(keyBuffer, ","))
	for _, section := range k.Content.SectionStrings() {
		if strings.HasPrefix(section, keys) {
			counter++
		}
	}

	return counter > 1
}

func (k *Keymap) Keys() func(yield func(*Key) bool) {
	var rowName string

	return func(yield func(*Key) bool) {
		for _, s := range k.Content.Sections() {
			if s.Name() == ini.DefaultSection {
				continue
			}

			if strings.HasPrefix(s.Name(), "--") {
				rowName = strings.Trim(s.Name(), "-")
				continue
			}

			key := k.NewKey(s.Name())

			if key == nil {
				continue
			}

			key.Row = rowName

			if !yield(key) {
				return
			}
		}
	}
}

func (k *Keymap) StoreKeyboard(path string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.New("failed to get current directory")
	}

	// Not using system temp dir because rename across filesystems isn't supported
	// and /tmp is probably on a separate partition.
	tempFile, err := os.CreateTemp(cwd, "keys-temp*.ini")
	if err != nil {
		return errors.New("failed to create temporary file")
	}
	defer func() {
		removeErr := os.Remove(tempFile.Name())
		if removeErr != nil {
			err = removeErr
		}
	}()

	k.Content.Section("").Key("keyboard").SetValue(path)

	err = k.Content.SaveTo(tempFile.Name())
	if err != nil {
		return fmt.Errorf("could not save file %s: %w", tempFile.Name(), err)
	}

	if err := os.Rename(tempFile.Name(), k.Filename); err != nil {
		return fmt.Errorf("could not open file %s: %w", tempFile.Name(), err)
	}

	return nil
}
