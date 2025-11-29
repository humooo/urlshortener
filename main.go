package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

type URLShortener struct {
	mu   sync.Mutex
	urls map[string]string
	keys map[string]string
}

var shortener = URLShortener{
	urls: make(map[string]string),
	keys: make(map[string]string),
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func main() {

	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/go/", handleRedirect)

	addr := ":" + *port
	log.Printf("Starting server on port %s...", *port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}

}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/go/")
	shortener.mu.Lock()
	url, ok := shortener.urls[key]
	shortener.mu.Unlock()
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)

}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	shortener.mu.Lock()
	defer shortener.mu.Unlock()

	var key string
	var ok bool

	if key, ok = shortener.keys[req.URL]; ok {
		resp := ShortenResponse{
			URL: req.URL,
			Key: key,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}

	for {
		key = generateKey()
		if _, ok = shortener.urls[key]; !ok {
			break
		}
	}

	shortener.urls[key] = req.URL
	shortener.keys[req.URL] = key

	resp := ShortenResponse{
		URL: req.URL,
		Key: key,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func generateKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
