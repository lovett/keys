package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Keymap struct {
	Filename string
	Content  *ini.File
}

func NewKeymap(filename string) *Keymap {
	c := Keymap{
		Filename: filename,
	}

	c.Parse()

	return &c
}

func (k *Keymap) Reload() {
	k.Parse()
}

func (c *Keymap) Parse() error {
	var err error

	options := ini.LoadOptions{
		SkipUnrecognizableLines: true,
	}

	if _, statErr := os.Stat(c.Filename); os.IsNotExist(statErr) {
		skeleton, err := ReadAsset("assets/skeleton.ini")
		if err != nil {
			return err
		}

		c.Content, err = ini.LoadSources(options, skeleton.Bytes)
	} else {
		c.Content, err = ini.LoadSources(options, c.Filename)
	}

	if err != nil {
		return err
	}

	// Usage is read-only, so (maybe?) speed up read operations.
	//
	// See https://ini.unknwon.io/docs/faqs
	c.Content.BlockMode = false

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
	sectionName := strings.Replace(keyName, "KEY_", "", 1)
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
	sectionName := k.KeyNameToSectionName(keyName)
	fmt.Println(sectionName)
	section, err := k.Content.GetSection(sectionName)

	if err != nil {
		return nil
	}

	return NewKeyFromSection(section)

}

func (c *Keymap) Keys() func(yield func(*Key) bool) {
	return func(yield func(*Key) bool) {
		for _, s := range c.Content.Sections() {
			if s.Name() == ini.DefaultSection {
				continue
			}

			k := NewKeyFromSection(s)

			if k == nil {
				continue
			}

			if !yield(k) {
				return
			}
		}
	}
}
