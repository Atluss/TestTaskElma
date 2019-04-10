package web_serve

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key         = []byte("super-secret-key")
	store       = sessions.NewCookieStore(key)
	sessionName = "session_a"
)

type HeadRequest interface {
	Request() // execute FetchTask
}

type RequestHttp struct {
	HeadRequest
	w *http.ResponseWriter
	r *http.Request
}

func setDefaultHeadersHttp(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
}
