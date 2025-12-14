package server

import (
	"fmt"
	"io"
	"keys/internal/asset"
	"keys/internal/config"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func failIfServerError(t *testing.T, rr *httptest.ResponseRecorder) {
	if rr.Code == http.StatusInternalServerError {
		t.Fatal(rr.Body.String())
	}
}

func resetLogger() {
	log.SetOutput(os.Stdout)
}

func tempFile(t *testing.T) *os.File {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tempFile, err := os.CreateTemp(wd, "keys-temp*.ini")
	if err != nil {
		t.Fatal(err)
	}

	return tempFile
}

func serverFixture(t *testing.T, fixture string) Server {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	fixturePath := filepath.Join(wd, "../../testdata", fixture)
	cfg, err := config.NewConfig(fixturePath)
	if err != nil {
		t.Fatal(err)
	}

	return Server{":4004", cfg}
}

func TestVersionHandler(t *testing.T) {
	server := serverFixture(t, "key-multiple.ini")
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.versionHandler)
	handler.ServeHTTP(rr, req)
	failIfServerError(t, rr)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("expected text/plain, got %s", contentType)
	}
}

func TestShellHandler(t *testing.T) {
	server := serverFixture(t, "key-multiple.ini")
	server.Config.PublicUrl = "https://example.com"

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.shellHandler)
	handler.ServeHTTP(rr, req)
	failIfServerError(t, rr)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("expected text/plain, got %s", contentType)
	}

	etag := rr.Header().Get("Etag")
	if len(etag) != 0 {
		t.Errorf("etag header should not have been set")
	}

	body := rr.Body.String()
	if !strings.Contains(body, server.Config.PublicUrl) {
		t.Errorf("response body did not contain publis url")
	}
}

func TestOpenApiHandler(t *testing.T) {
	server := serverFixture(t, "key-multiple.ini")
	server.Config.PublicUrl = "https://example.com"

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.openapiHandler)
	handler.ServeHTTP(rr, req)
	failIfServerError(t, rr)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	accessControl := rr.Header().Get("Access-Control-Allow-Origin")
	if accessControl != "*" {
		t.Errorf("expected Access-Control-Allow-Origin to be *, got %s", accessControl)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/yaml" {
		t.Errorf("expected application/yaml, got %s", contentType)
	}

	tests := []struct {
		name   string
		search string
	}{
		{name: "public url", search: fmt.Sprintf("url: \"%s\"", server.Config.PublicUrl)},
		{name: "version path", search: "/version:"},
	}

	body := rr.Body.String()
	for _, tt := range tests {
		if !strings.Contains(body, tt.search) {
			t.Errorf("response body did not reflect %s", tt.name)
		}
	}

}

func TestAssetHandler(t *testing.T) {
	server := serverFixture(t, "key-multiple.ini")

	tests := []struct {
		path        string
		contentType string
		code        int
	}{
		{"/assets/favicon.svg", "image/svg+xml", 200},
		{"/assets/keys.css", "text/css", 200},
		{"/assets/keys.js", "application/javascript", 200},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.assetHandler)
		handler.ServeHTTP(rr, req)
		failIfServerError(t, rr)

		if rr.Code != tt.code {
			t.Errorf("expected %d, got %d", tt.code, rr.Code)
		}

		contentType := rr.Header().Get("Content-Type")
		if contentType != tt.contentType {
			t.Errorf("expected %s, got %s", tt.contentType, contentType)
		}

		etag := rr.Header().Get("Etag")
		if len(etag) != 32 {
			t.Errorf("etag header was not a 32-char string. got %s", etag)
		}
	}
}

func TestAssetHandlerBogusHash(t *testing.T) {
	server := serverFixture(t, "key-multiple.ini")
	req := httptest.NewRequest("GET", "/assets/keys.js", nil)
	req.Header.Set("If-None-Match", "my-bogus-hash")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.assetHandler)
	handler.ServeHTTP(rr, req)
	failIfServerError(t, rr)

	if rr.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	etag := rr.Header().Get("Etag")
	if len(etag) != 32 {
		t.Errorf("etag value was not a 32-char string: '%s'", etag)
	}
}

func TestAssetHandlerIfNoneMatch(t *testing.T) {
	server := serverFixture(t, "key-multiple.ini")

	tests := []struct {
		path string
	}{
		{"/assets/favicon.svg"},
		{"/assets/keys.css"},
		{"/assets/keys.js"},
	}

	for _, tt := range tests {
		asset, err := asset.Read(strings.TrimPrefix(tt.path, "/"))
		if err != nil {
			t.Fatal(err)
		}
		if asset.Hash == "" {
			t.Fatalf("no hash for %s with mime type %s", asset.Path, asset.MimeType)
		}

		req := httptest.NewRequest("GET", tt.path, nil)
		req.Header.Set("If-None-Match", asset.Hash)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.assetHandler)
		handler.ServeHTTP(rr, req)
		failIfServerError(t, rr)

		if rr.Code != http.StatusNotModified {
			t.Errorf("got %d instead of 304", rr.Code)
		}
	}
}

func TestTriggerHandler(t *testing.T) {
	t.Cleanup(resetLogger)
	log.SetOutput(io.Discard)

	server := serverFixture(t, "key-multiple.ini")

	tests := []struct {
		name string
		key  string
		code int
	}{
		{"invalid", "", http.StatusNotFound},
		{"test", "", http.StatusOK},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("POST", "/", nil)
		req.SetPathValue("name", tt.name)
		req.SetPathValue("key", tt.key)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.triggerHandler)
		handler.ServeHTTP(rr, req)
		failIfServerError(t, rr)

		if rr.Code != tt.code {
			t.Errorf("expected %d, got %d", tt.code, rr.Code)
		}

		if rr.Header().Get("X-Keys-Locked") != "" {
			t.Error("X-Keys-Locked header found on non-lock request")
		}
	}
}

func TestTriggerLock(t *testing.T) {
	t.Cleanup(resetLogger)
	log.SetOutput(io.Discard)

	server := serverFixture(t, "key-roll-lock.ini")
	req := httptest.NewRequest("POST", "/", nil)
	req.SetPathValue("name", "test")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.triggerHandler)
	handler.ServeHTTP(rr, req)
	failIfServerError(t, rr)

	if server.Config.KeyboardLocked != true {
		t.Error("Keyboard was not locked")
	}

	if rr.Header().Get("X-Keys-Locked") != "1" {
		t.Error("X-Keys-Locked header not 1")
	}

	body := rr.Body.String()
	if body != "Keyboard locked" {
		t.Errorf("Unexpected body from lock request: '%s'", body)
	}

	req2 := httptest.NewRequest("POST", "/", nil)
	req2.SetPathValue("name", "test")
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	failIfServerError(t, rr)

	if server.Config.KeyboardLocked != false {
		t.Error("Keyboard was not unlocked")
	}

	if rr2.Header().Get("X-Keys-Locked") != "0" {
		t.Error("X-Keys-Locked header not 0")
	}

	body = rr2.Body.String()
	if body != "Keyboard unlocked" {
		t.Errorf("Unexpected body from lock request: '%s'", body)
	}

}

func TestTriggerToggle(t *testing.T) {
	t.Cleanup(resetLogger)
	log.SetOutput(io.Discard)

	server := serverFixture(t, "key-roll.ini")

	tests := []struct {
		responseBody string
		resultState  string
	}{
		{"hello\n", "state2"},
		{"hello 2\n", "state3"},
		{"hello 3\n", "state1"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("POST", "/", nil)
		req.SetPathValue("name", "test")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.triggerHandler)
		handler.ServeHTTP(rr, req)
		failIfServerError(t, rr)

		state := rr.Header().Get("X-Keys-State")
		if state != tt.resultState {
			t.Errorf("wanted state after first request to be '%s', got '%s'", tt.resultState, state)
		}

		body := rr.Body.String()
		if body != tt.responseBody {
			t.Errorf("wanted response body of first request to be '%s', got '%s'", tt.responseBody, body)
		}
	}
}

func TestKeymapHandler(t *testing.T) {
	server := serverFixture(t, "key-multiple.ini")

	tests := []struct {
		contentType string
		bodyMatch   string
		status      int
	}{
		{"text/plain", "test2 (w)\n  echo hello world", 200},
		{"text/html", "href=\"/trigger/test2\"", 200},
		{"", "href=\"/trigger/test2\"", 200},
		{"garbage", "Unsupported", 406},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/", nil)
		if tt.contentType != "" {
			req.Header.Set("Accept", tt.contentType)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.keymapHandler)
		handler.ServeHTTP(rr, req)
		failIfServerError(t, rr)

		if rr.Code != tt.status {
			t.Errorf("expected %d for content type '%s', got %d", tt.status, tt.contentType, rr.Code)
		}

		contentType := rr.Header().Get("Content-Type")
		wantedContentType := tt.contentType
		if tt.status == http.StatusNotAcceptable {
			wantedContentType = "text/plain"
		}

		if tt.status == http.StatusOK {
			if tt.contentType != tt.contentType {
				t.Errorf("expected %s, got %s", wantedContentType, contentType)
			}
		}

		if tt.bodyMatch != "" {
			body := rr.Body.String()
			if !strings.Contains(body, tt.bodyMatch) {
				t.Errorf("response body did not contain '%s': %s", tt.bodyMatch, body)
			}
		}
	}
}

func TestEditHandler(t *testing.T) {
	server := serverFixture(t, "key-multiple.ini")

	tests := []struct {
		contentType string
		bodyMatch   string
		status      int
	}{
		{"", string(server.Config.Keymap.Raw()), 200},
		{"garbage", "Unsupported", 406},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/", nil)
		if tt.contentType != "" {
			req.Header.Set("Accept", tt.contentType)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.editHandler)
		handler.ServeHTTP(rr, req)
		failIfServerError(t, rr)

		if rr.Code != tt.status {
			t.Errorf("expected %d for content type '%s', got %d", tt.status, tt.contentType, rr.Code)
		}
		failIfServerError(t, rr)
	}
}

func TestSaveHandler(t *testing.T) {
	tmpFile := tempFile(t)

	t.Cleanup(func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Fatal(err)
		}
	})

	configBody := `[temp]
command = echo temp
physical_key = t
`
	configBody2 := `
sound = off

[tempedit]
command = echo temp2
physical_key = t2
`

	_, err := tmpFile.WriteString(configBody)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.NewConfig(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	server := Server{":4004", cfg}

	form := url.Values{}
	form.Set("content", configBody2)

	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.saveHandler)
	handler.ServeHTTP(rr, req)
	failIfServerError(t, rr)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected %d, got %d", http.StatusSeeOther, rr.Code)
	}

	if key := server.Config.Keymap.FindKey("temp"); key != nil {
		t.Errorf("config was not reloaded after edit (old key found)")
	} else if key := server.Config.Keymap.FindKey("tempedit"); key == nil {
		t.Errorf("config was not reloaded after edit (new key not found)")
	} else if server.Config.Keymap.SoundAllowed == true {
		t.Errorf("config was not reloaded after edit (sound on)")
	} else if key := server.Config.Keymap.FindKey("t2"); key == nil {
		t.Errorf("config was not reloaded after edit (new physical key not found)")
	}
}
