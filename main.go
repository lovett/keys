package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const APP_VERSION = "2025.09.27"

func help() {
	fmt.Fprint(os.Stderr, "Relay keyboard input to commands or services.\n\n")

	fmt.Fprint(os.Stderr, "Options\n")
	flag.PrintDefaults()
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)

	asset, err := ReadAsset(strings.TrimPrefix(r.RequestURI, "/"))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Header.Get("If-None-Match") == asset.Hash {
		w.WriteHeader(http.StatusNotModified)
	} else {
		w.Header().Set("Content-Type", asset.MimeType)
		w.Header().Set("ETag", asset.Hash)
		w.Write(asset.Bytes)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)

	config := NewConfig()
	config.Read()

	html, err := config.RenderEdit()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(html)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)

	config := NewConfig()

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form.", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Form values: %v\n", r.Form)

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

	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(r.Form.Get("content"))); err != nil {
		http.Error(w, "Failed to write to temporary file.", http.StatusInternalServerError)
		return
	}

	if err := tempFile.Close(); err != nil {
		http.Error(w, "Failed to close temp file.", http.StatusInternalServerError)
		return
	}

	if err := os.Rename(tempFile.Name(), config.Filename); err != nil {
		fmt.Println(fmt.Errorf("could not open file %q: %w", tempFile.Name(), err))
		http.Error(w, "Failed to move temporary file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)

	config := NewConfig()
	config.Parse()

	html, err := config.RenderKeyboard()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Write(html)
	}
}

func triggerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)

	config := NewConfig()
	config.Parse()

	stdout, err := config.Fire(r.PathValue("trigger"))

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(err.Error()))
		return
	}

	if len(stdout) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(stdout)
}

func main() {
	version := flag.Bool("version", false, "Application version")

	flag.Usage = help
	flag.Parse()
	if *version {
		fmt.Println(APP_VERSION)
		os.Exit(0)
	}

	if err := HashAssets(); err != nil {
		log.Fatalf("Error during asset hashing: %s", err.Error())
		return
	}

	http.HandleFunc("GET /{$}", dashboardHandler)
	http.HandleFunc("GET /assets/favicon.svg", assetHandler)
	http.HandleFunc("GET /assets/keys.css", assetHandler)
	http.HandleFunc("GET /assets/keys.js", assetHandler)
	http.HandleFunc("GET /edit", editHandler)
	http.HandleFunc("POST /edit", saveHandler)
	http.HandleFunc("POST /{trigger}", triggerHandler)

	host := os.Getenv("KEYS_HOST")

	port := os.Getenv("KEYS_PORT")
	if port == "" {
		port = "4004"
	}

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Serving on %s", address)
	log.Fatal(http.ListenAndServe(address, nil))

}
