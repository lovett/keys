package asset

import (
	"crypto/md5"
	"embed"
	"fmt"
	"log"
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

func ReadVersion() []byte {
	b, err := AssetFS.ReadFile("assets/version.txt")

	if err != nil {
		return []byte("dev")
	}

	return b
}

func ReadAsset(path string) (*Asset, error) {
	b, err := AssetFS.ReadFile(path)

	if err != nil {
		return nil, err
	}

	mime := mimeType(path)

	hash := ""
	if mime != "" {
		hash = assetHash(path)
	}

	asset := Asset{
		Path:     path,
		Bytes:    b,
		Hash:     hash,
		MimeType: mime,
	}

	return &asset, nil
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
		return ""
	}
}

func assetHash(path string) string {
	if hash, ok := hashCache[path]; ok {
		return hash
	}

	b, err := AssetFS.ReadFile(path)

	if err != nil {
		log.Printf("Error during asset hashing: %s", err)
		hashCache[path] = ""
	} else {
		hashCache[path] = fmt.Sprintf("%x", md5.Sum(b))
	}

	return hashCache[path]
}
