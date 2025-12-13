package asset

import (
	"testing"
)

func clearCache() {
	clear(hashCache)
}

func TestAssetRead(t *testing.T) {
	t.Cleanup(clearCache)

	tests := []struct {
		path string
		mime string
	}{
		{path: "assets/keys.css", mime: "text/css"},
		{path: "assets/keys.js", mime: "application/javascript"},
		{path: "assets/favicon.svg", mime: "image/svg+xml"},
		{path: "assets/openapi.yaml", mime: ""},
	}

	for _, tt := range tests {
		asset, err := Read(tt.path)
		if err != nil {
			t.Fatalf("Could not read %s: %v", tt.path, err)
		}

		if asset.MimeType != tt.mime {
			t.Errorf("Expected %s for %s, got %s", tt.mime, tt.path, asset.MimeType)
		}

		if asset.MimeType == "" && asset.Hash != "" {
			t.Errorf("Asset without mime should not be hashed. Got %s for %s", asset.Hash, tt.path)
		}

		if asset.MimeType != "" && asset.Hash == "" {
			t.Errorf("%s should have been hashed because it has a mime type, but wasn't", tt.path)
		}

		if asset.Path != tt.path {
			t.Errorf("Path for %s is unexpectedly %s", tt.path, asset.Path)
		}
	}
}

func TestAssetExistence(t *testing.T) {
	t.Cleanup(clearCache)

	asset, _ := Read("does-not-exist")
	if asset != nil {
		t.Errorf("Nonexistant asset path should have been rejected")
	}
}

func TestVersionDefault(t *testing.T) {
	t.Cleanup(clearCache)

	b := ReadVersion()

	if string(b) != "unknown" {
		t.Errorf("Unexpected default value for version asset: %v", b)
	}
}

func TestAssetHashMatching(t *testing.T) {
	t.Cleanup(clearCache)

	tests := []struct {
		hash  string
		match bool
	}{
		{hash: "my-hash", match: true},
		{hash: "", match: false},
	}

	a, err := Read("assets/keyboard.html")
	if err != nil {
		t.Fatalf("Could not read asset: %v", err)
	}

	for _, tt := range tests {
		a.Hash = tt.hash
		result := a.HashMatch(tt.hash)
		if result != tt.match {
			t.Fatalf("Hash matching failure: got %t for %s", result, tt.hash)
		}
	}
}

func TestAssetHashCaching(t *testing.T) {
	t.Cleanup(clearCache)

	path := "assets/keys.css"
	a1, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if a1.Hash == "" {
		t.Fatal("Asset hash wasn't set after first read")
	}

	a2, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}

	if a2.Hash == "" {
		t.Errorf("Asset hash wasn't set after second read")
	}

	if a1.Hash != a2.Hash {
		t.Errorf("Asset hash mismatch between first and second read")
	}
}
