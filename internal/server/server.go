package server

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"keys/internal/asset"
	"keys/internal/config"
	"keys/internal/keymap"
	"keys/internal/sound"
	"log"
	"net/http"
	"strings"
	texttemplate "text/template"
	"time"
)

type Server struct {
	ServerAddress string
	Config        *config.Config
}

func Serve(cfg *config.Config, port int) {
	s := Server{
		ServerAddress: fmt.Sprintf(":%d", port),
		Config:        cfg,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", s.keymapHandler)
	mux.HandleFunc("GET /assets/favicon.svg", s.assetHandler)
	mux.HandleFunc("GET /assets/keys.css", s.assetHandler)
	mux.HandleFunc("GET /assets/keys.js", s.assetHandler)
	mux.HandleFunc("GET /edit", s.editHandler)
	mux.HandleFunc("GET /openapi.yaml", s.openapiHandler)
	mux.HandleFunc("GET /version", s.versionHandler)
	mux.HandleFunc("POST /edit", s.saveHandler)
	mux.HandleFunc("POST /trigger/{name}", s.triggerHandler)
	mux.HandleFunc("GET /util/sh", s.shellHandler)
	log.Printf("Serving on %s and available from %s", s.ServerAddress, cfg.PublicUrl)
	log.Printf("Config file is %s", cfg.Keymap.Filename)

	server := &http.Server{
		Addr:         s.ServerAddress,
		Handler:      requestLogger(serverHeaders(mux, s.Config)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func serverHeaders(next http.Handler, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") != "" && r.Header.Get("X-Forwarded-Host") != "" {
			forwardedPublicUrl := fmt.Sprintf("%s://%s", r.Header.Get("X-Forwarded-Proto"), r.Header.Get("X-Forwarded-Host"))
			if forwardedPublicUrl != config.PublicUrl {
				config.PublicUrl = forwardedPublicUrl
			}
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Server", "keys")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) acceptableRequest(w http.ResponseWriter, r *http.Request, acceptableContentTypes []string) bool {
	accept := r.Header.Get("Accept")
	if accept == "" {
		accept = "text/html"
	}

	for _, contentType := range acceptableContentTypes {
		if strings.Contains(accept, contentType) {
			return true
		}
	}

	http.Error(w, "Unsupported content type.", http.StatusNotAcceptable)
	return false
}

func (s *Server) keymapHandler(w http.ResponseWriter, r *http.Request) {
	if !s.acceptableRequest(w, r, []string{"text/html", "text/plain"}) {
		return
	}

	if r.Header.Get("Accept") == "text/plain" {
		s.keymapTextWriter(w, r)
	}
	s.keymapHtmlWriter(w)
}

func (s *Server) keymapTextWriter(w http.ResponseWriter, r *http.Request) {
	var output bytes.Buffer

	query := r.URL.Query()
	label := strings.ToLower(query.Get("label"))
	command := strings.ToLower(query.Get("command"))
	keyboardKey := strings.ToLower(query.Get("key"))

	funcMap := texttemplate.FuncMap{
		"queryMatch": func(k keymap.Key) bool {
			if label != "" && !k.MatchesName(label) {
				return false
			}

			if command != "" && !k.MatchesCommand(command) {
				return false
			}

			if keyboardKey != "" && !k.MatchesPhysicalKey(keyboardKey) {
				return false
			}

			return true
		},
	}

	tmpl := texttemplate.New("keyboard.txt").Funcs(funcMap)
	tmpl, err := tmpl.ParseFS(asset.AssetFS, "assets/keyboard.txt")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("unable to parse text template: %v", err)
		return
	}

	if err := tmpl.ExecuteTemplate(&output, "keyboard.txt", s.Config); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("unable to write error response body: %v", err)
		}
	} else {
		w.Header().Set("Content-Type", "text/plain")
		if _, err = w.Write(output.Bytes()); err != nil {
			log.Fatalf("unable to write response body: %v", err)
		}
	}
}

func (s *Server) keymapHtmlWriter(w http.ResponseWriter) {
	templates := htmltemplate.Must(htmltemplate.ParseFS(asset.AssetFS, "assets/layout.html", "assets/keyboard.html"))

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

func (s *Server) shellHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := texttemplate.New("keys.sh")
	tmpl, err := tmpl.ParseFS(asset.AssetFS, "assets/keys.sh")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("unable to parse shell template: %v", err)
		return
	}

	var output bytes.Buffer
	if err := tmpl.ExecuteTemplate(&output, "keys.sh", s.Config); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("unable to write error response body: %v", err)
		}
	} else {
		w.Header().Set("Content-Type", "text/plain")
		if _, err = w.Write(output.Bytes()); err != nil {
			log.Fatalf("unable to write response body: %v", err)
		}
	}
}

func (s *Server) assetHandler(w http.ResponseWriter, r *http.Request) {
	asset, err := asset.Read(strings.TrimPrefix(r.RequestURI, "/"))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	wantedHash := r.Header.Get("If-None-Match")

	if asset.HashMatch(wantedHash) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", asset.MimeType)
	w.Header().Set("ETag", asset.Hash)
	if _, err := w.Write(asset.Bytes); err != nil {
		log.Fatalf("unable to write asset body: %v", err)
	}
}

func (s *Server) editHandler(w http.ResponseWriter, r *http.Request) {
	if !s.acceptableRequest(w, r, []string{"text/html"}) {
		return
	}

	templates := htmltemplate.Must(htmltemplate.ParseFS(asset.AssetFS, "assets/layout.html", "assets/editor.html"))

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
	err := r.ParseForm()
	if err != nil {
		wrappedError := fmt.Errorf("error during form parsing: %w", err)
		http.Error(w, wrappedError.Error(), http.StatusInternalServerError)
		return
	}

	content := []byte(r.Form.Get("content"))
	err = s.Config.Keymap.Replace(content)
	if err != nil {
		wrappedError := fmt.Errorf("error during save: %w", err)
		http.Error(w, wrappedError.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) maybePlaySound(name sound.Name) {
	if s.Config.Keymap.SoundAllowed {
		return
	}

	if err := sound.Play(name); err != nil {
		log.Println(err)
	}
}

func (s *Server) triggerHandler(w http.ResponseWriter, r *http.Request) {
	key := s.Config.Keymap.FindKey(r.PathValue("name"))
	if key == nil {
		s.maybePlaySound(sound.Error)
		http.NotFound(w, r)
		return
	}

	var stdout []byte
	var err error
	switch key.CurrentCommand() {
	case "lock":
		s.maybePlaySound(sound.Lock)
		s.Config.KeyboardLocked = true
		key.Toggle()
		stdout = []byte("Keyboard locked")
		w.Header().Set("X-Keys-Locked", "1")
	case "unlock":
		s.maybePlaySound(sound.Unlock)
		s.Config.KeyboardLocked = false
		key.Toggle()
		stdout = []byte("Keyboard unlocked")
		w.Header().Set("X-Keys-Locked", "0")
	default:
		s.maybePlaySound(sound.Tap)
		stdout, err = key.RunCommand()
		if err != nil {
			s.maybePlaySound(sound.Error)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if key.Confirmation {
			s.maybePlaySound(sound.Confirmation)
		}
	}

	if key.CanToggle() {
		w.Header().Set("X-Keys-State", key.State())
	}

	if len(stdout) == 0 || !key.ShowOutput {
		w.WriteHeader(http.StatusNoContent)
		return
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
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write(asset.ReadVersion()); err != nil {
		log.Fatalf("unable to write version response body: %v", err)
	}
}

func (s *Server) openapiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	tmpl := texttemplate.New("openapi.yaml")
	tmpl, err := tmpl.ParseFS(asset.AssetFS, "assets/openapi.yaml")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("unable to parse template: %v", err)
		return
	}

	var output bytes.Buffer
	if err := tmpl.ExecuteTemplate(&output, "openapi.yaml", s.Config); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			log.Fatalf("unable to write error response body: %v", err)
		}
	} else {
		w.Header().Set("Content-Type", "application/yaml")
		if _, err = w.Write(output.Bytes()); err != nil {
			log.Fatalf("unable to write response body: %v", err)
		}
	}
}
