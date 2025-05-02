package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/scenery/mediax/config"
)

const AuthCookieName = "session"

type SessionData struct {
	// ID            int
	// CreatedAt     time.Time
	Authenticated bool
}

// In-memory session storage
var sessions = make(map[string]SessionData)
var sessionsMutex sync.Mutex

func CreateSession(w http.ResponseWriter) error {
	token, err := generateSessionToken()
	if err != nil {
		return fmt.Errorf("failed to generate session token: %w", err)
	}

	sessionsMutex.Lock()
	sessions[token] = SessionData{
		Authenticated: true,
	}
	sessionsMutex.Unlock()

	cookie := &http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(config.App.SessionTimeout),
		HttpOnly: true,
		Secure:   config.App.Server.UseHTTPS,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	return nil
}

func DeleteSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(AuthCookieName)
	if err == nil {
		token := cookie.Value
		sessionsMutex.Lock()
		delete(sessions, token)
		sessionsMutex.Unlock()
	} else {
		log.Println("No auth cookie found in request to delete session.")
	}

	clearClientCookie(w)
}

func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie(AuthCookieName)
	if err != nil {
		return false
	}

	token := cookie.Value

	sessionsMutex.Lock()
	session, ok := sessions[token]
	sessionsMutex.Unlock()
	if !ok {
		return false
	}

	return session.Authenticated
}

func generateSessionToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func clearClientCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     AuthCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   config.App.Server.UseHTTPS,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}
