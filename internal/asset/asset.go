package asset

import (
	"crypto/md5"
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed assets/*
var AssetFS embed.FS

type Asset struct {
	Path     string
	MimeType string
	Bytes    []byte
	Hash     string
}

var hashCache = make(map[string]string)

func Read(path string) (*Asset, error) {
	b, err := AssetFS.ReadFile(path)

	if err != nil {
		return nil, err
	}

	asset := Asset{
		Path:  path,
		Bytes: b,
	}

	switch filepath.Ext(asset.Path) {
	case ".css":
		asset.MimeType = "text/css"
	case ".js":
		asset.MimeType = "application/javascript"
	case ".svg":
		asset.MimeType = "image/svg+xml"
	}

	hash, found := hashCache[asset.Path]

	if asset.MimeType != "" && !found {
		hash = fmt.Sprintf("%x", md5.Sum(asset.Bytes))
		hashCache[asset.Path] = hash
	}

	asset.Hash = hash

	return &asset, nil
}

func ReadVersion() []byte {
	asset, err := Read("assets/version.txt")

	if err != nil {
		return []byte("unknown")
	}

	return asset.Bytes
}
