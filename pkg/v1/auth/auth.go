// Package auth session realisation to auth
package auth

import (
	"github.com/Atluss/TestTaskElma/pkg/v1"
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	Key         = []byte("super-secret-key")
	Store       = sessions.NewCookieStore(Key)
	SessionName = "session_a"
)

// GetSession check session
func GetSession(r *http.Request) *sessions.Session {
	session, _ := Store.Get(r, SessionName)
	return session
}

// CheckAuth check session
func CheckAuth(session *sessions.Session) bool {
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false
	}
	return true
}

// SaveSessionLogin login by session
func SaveSessionLogin(session *sessions.Session, w http.ResponseWriter, r *http.Request) {
	session.Values["authenticated"] = true
	v1.LogOnError(session.Save(r, w), "warning")
}

// SaveSessionLogout logout by session
func SaveSessionLogout(session *sessions.Session, w http.ResponseWriter, r *http.Request) {
	session.Values["authenticated"] = false
	v1.LogOnError(session.Save(r, w), "warning")
}
