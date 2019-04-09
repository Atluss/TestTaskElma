package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)

	sessionName = "session_a"
)

func main() {

	route := mux.NewRouter().StrictSlash(true)

	route.HandleFunc("/", rootHttp)
	route.HandleFunc("/login", login)
	route.HandleFunc("/logout", logout)

	http.ListenAndServe(":8080", route)

}

func rootHttp(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, sessionName)
	w.Header().Set("Content-Type", "text/html")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		w.WriteHeader(http.StatusForbidden)
		http.ServeFile(w, r, "public/login.html")
	} else {
		http.ServeFile(w, r, "public/list.html")
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}
