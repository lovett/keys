package server

import (
	"bytes"
	"fmt"
	"html/template"
	"keys/internal/asset"
	"keys/internal/config"
	"keys/internal/sound"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	Config *config.Config
}

func StartServer(config *config.Config) {
	s := Server{
		Config: config,
	}

	if err := asset.HashAssets(); err != nil {
		log.Fatalf("Error during asset hashing: %s", err.Error())
	}

	http.HandleFunc("GET /{$}", s.dashboardHandler)
	http.HandleFunc("GET /assets/favicon.svg", s.assetHandler)
	http.HandleFunc("GET /assets/keys.css", s.assetHandler)
	http.HandleFunc("GET /assets/keys.js", s.assetHandler)
	http.HandleFunc("GET /edit", s.editHandler)
	http.HandleFunc("GET /version", s.versionHandler)
	http.HandleFunc("POST /edit", s.saveHandler)
	http.HandleFunc("POST /trigger/{key}", s.triggerHandler)
	log.Printf("Serving on %s and available from %s", config.ServerAddress, config.PublicUrl)
	log.Printf("Config file is %s", config.Keymap.Filename)
	log.Fatal(http.ListenAndServe(config.ServerAddress, nil))
}

func (s *Server) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)

	templates := template.Must(template.ParseFS(asset.AssetFS, "assets/layout.html", "assets/keyboard.html"))

	var output bytes.Buffer
	if err := templates.ExecuteTemplate(&output, "layout.html", s.Config); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("unable to write error response body: %v", err)
		}
	} else {
		w.Header().Set("Content-Type", "text/html")
		if _, err = w.Write(output.Bytes()); err != nil {
			log.Fatalf("unable to write response body: %v", err)
		}
	}
}

func (s *Server) assetHandler(w http.ResponseWriter, r *http.Request) {
	asset, err := asset.ReadAsset(strings.TrimPrefix(r.RequestURI, "/"))

	if err != nil {
		s.logRequestWithStatus(r, http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Header.Get("If-None-Match") == asset.Hash {
		s.logRequestWithStatus(r, http.StatusNotModified)
		w.WriteHeader(http.StatusNotModified)
	} else {
		s.logRequestWithStatus(r, http.StatusOK)
		w.Header().Set("Content-Type", asset.MimeType)
		w.Header().Set("ETag", asset.Hash)
		if _, err := w.Write(asset.Bytes); err != nil {
			log.Fatalf("unable to write asset body: %v", err)
		}
	}
}

func (s *Server) editHandler(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)

	templates := template.Must(template.ParseFS(asset.AssetFS, "assets/layout.html", "assets/editor.html"))

	var output bytes.Buffer

	if err := templates.ExecuteTemplate(&output, "layout.html", s.Config.Keymap); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("unable to write edit response error body: %v", err)
		}
	} else {
		if _, err := w.Write(output.Bytes()); err != nil {
			log.Fatalf("unable to write edit response body: %v", err)
		}

	}
}

func (s *Server) saveHandler(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form.", http.StatusInternalServerError)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Failed to get current directory", http.StatusInternalServerError)
		return
	}

	// Not using system temp dir because rename across filesystems isn't supported
	// and /tmp is probably on a separate partition.
	tempFile, err := os.CreateTemp(cwd, "keys-temp*.ini")

	if err != nil {
		http.Error(w, "Failed to create temporary file.", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			if os.IsNotExist(err) {
				return
			}
			log.Fatalf("unable to remove tempfile: %v", err)
		}
	}()

	if _, err := tempFile.Write([]byte(r.Form.Get("content"))); err != nil {
		http.Error(w, "Failed to write to temporary file.", http.StatusInternalServerError)
		return
	}

	if err := tempFile.Close(); err != nil {
		http.Error(w, "Failed to close temp file.", http.StatusInternalServerError)
		return
	}

	if err := os.Rename(tempFile.Name(), s.Config.Keymap.Filename); err != nil {
		fmt.Println(fmt.Errorf("could not open file %q: %w", tempFile.Name(), err))
		http.Error(w, "Failed to move temporary file", http.StatusInternalServerError)
		return
	}

	err = s.Config.Keymap.Reload()
	if err != nil {
		log.Fatalf("Error during reload: %v", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) triggerHandler(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)

	key := s.Config.Keymap.NewKey(r.PathValue("key"))

	if key == nil {
		s.sendError(w, "Invalid key")
		return
	}

	var stdout []byte
	var err error
	switch key.CurrentCommand() {
	case "lock":
		s.Config.KeyboardLocked = true
		key.UpdateCommandIndex()
		stdout = []byte("Keyboard locked")
		w.Header().Set("X-Keys-Locked", "1")
	case "unlock":
		s.Config.KeyboardLocked = false
		key.UpdateCommandIndex()
		stdout = []byte("Keyboard unlocked")
		w.Header().Set("X-Keys-Locked", "0")
	default:
		stdout, err = key.RunCommand()
		if err != nil {
			sound.PlayErrorSound(s.Config)
			s.sendError(w, err.Error())
			return
		}
	}

	if len(stdout) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	if key.Toggle {
		sound.PlayToggleSound(s.Config, key)
		w.Header().Set("X-Keys-State", key.CurrentState())
	} else {
		sound.PlayConfirmationSound(s.Config)
	}

	if bytes.ContainsRune(stdout, '<') && bytes.ContainsRune(stdout, '>') {
		w.Header().Set("Content-Type", "text/html")
	} else {
		w.Header().Set("Content-Type", "text/plain")
	}

	if _, err := w.Write(stdout); err != nil {
		log.Fatalf("unable to write stdout response body: %v", err)
	}
}

func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(s.Config.AppVersion)); err != nil {
		log.Fatalf("unable to write version response body: %v", err)
	}
}

func (s *Server) sendError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(message)); err != nil {
		log.Fatalf("unable to write error response body: %v", err)
	}
}

func (s *Server) logRequest(r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)
}

func (s *Server) logRequestWithStatus(r *http.Request, status int) {
	log.Printf("%s %s -> %d", r.Method, r.RequestURI, status)
}
