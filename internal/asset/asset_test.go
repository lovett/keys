package asset

import (
	"testing"
)

func TestAssetRead(t *testing.T) {
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
		asset, err := ReadAsset(tt.path)
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
	asset, _ := ReadAsset("does-not-exist")
	if asset != nil {
		t.Errorf("Nonexistant asset path should have been rejected")
	}
}

func TestVersionDefault(t *testing.T) {
	b := ReadVersion()

	if string(b) != "dev" {
		t.Errorf("Unexpected default value for version asset: %v", b)
	}
}
