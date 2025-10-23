package main

import (
	"crypto/md5"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/vorbis"
)

//go:embed assets/*
var AssetFS embed.FS

type Asset struct {
	Path     string
	MimeType string
	Bytes    []byte
	Hash     string
}

var hashMap = make(map[string]string)

func HashAssets() error {
	return fs.WalkDir(AssetFS, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == "html" {
			return nil
		}

		if filepath.Ext(path) == "ogg" {
			return nil
		}

		if filepath.Ext(path) == "ini" {
			return nil
		}

		b, err := AssetFS.ReadFile(path)

		if err != nil {
			return err
		}

		hashMap[path] = fmt.Sprintf("%x", md5.Sum(b))
		return nil
	})
}

func ReadAsset(path string) (*Asset, error) {
	b, err := AssetFS.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return &Asset{
		Path:     path,
		Bytes:    b,
		Hash:     hashMap[path],
		MimeType: mimeType(path),
	}, nil
}

func SoundBuffer(path string) *Sound {
	b, err := AssetFS.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := vorbis.Decode(b)
	if err != nil {
		log.Fatal(err)
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	return &Sound{
		Path:   path,
		Format: format,
		Buffer: *buffer,
	}
}

func mimeType(path string) string {
	switch filepath.Ext(path) {
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
