package main

import (
	"bytes"
	"errors"
	"html/template"
	"os"

	"gopkg.in/ini.v1"
)

type Config struct {
	RawContent []byte
	Filename   string
	Content    *ini.File
}

func NewConfig() *Config {
	c := Config{
		Filename: "keys.ini",
	}

	return &c
}

func (c *Config) Parse() error {
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

func (c *Config) Read() error {
	var err error

	if _, statErr := os.Stat(c.Filename); os.IsNotExist(statErr) {
		asset, err := ReadAsset("assets/skeleton.ini")
		if err != nil {
			return err
		}

		c.RawContent = asset.Bytes
	} else {
		c.RawContent, err = os.ReadFile(c.Filename)
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Keys() func(yield func(*Key) bool) {
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

func (c *Config) Fire(trigger string) ([]byte, error) {
	section, err := c.Content.GetSection(trigger)

	if err != nil {
		return nil, errors.New("Unknown key")
	}

	k := NewKeyFromSection(section)

	if k == nil {
		return nil, errors.New("Invalid key")
	}

	return k.RunCommand()
}

func (c *Config) RenderKeyboard() ([]byte, error) {
	var templates = template.Must(template.ParseFS(AssetFS, "assets/layout.html", "assets/keyboard.html"))

	var output bytes.Buffer
	if err := templates.ExecuteTemplate(&output, "layout.html", c); err != nil {
		return []byte{}, err
	}

	return output.Bytes(), nil
}

func (c *Config) RenderEdit() ([]byte, error) {
	var templates = template.Must(template.ParseFS(AssetFS, "assets/layout.html", "assets/editor.html"))

	var output bytes.Buffer
	if err := templates.ExecuteTemplate(&output, "layout.html", c); err != nil {
		return []byte{}, err
	}

	return output.Bytes(), nil
}
