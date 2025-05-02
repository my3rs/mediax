package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/scenery/mediax/handlers"
)

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; form-action 'self'; img-src 'self' data: https:; script-src 'sha256-osjxnKEPL/pQJbFk1dKsF7PYFmTyMWGmVSiL9inhxJY=' 'unsafe-hashes'; style-src 'self' 'unsafe-inline';")
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsAuthenticated(r) {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	})
}

const BearerTokenPrefix = "Bearer "

func APIAuthMiddleware(expectedHashedKey string) func(http.Handler) http.Handler {
	expectedHashedKeyBytes, err := base64.StdEncoding.DecodeString(expectedHashedKey)
	if err != nil {
		log.Fatalf("Error: Invalid Base64 encoding for api_key '%s': %v", expectedHashedKey, err)
	}

	if expectedHashedKey == "" || len(expectedHashedKeyBytes) == 0 {
		log.Println("Info: API Key is not configured. All API requests will be denied with 403.")
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlers.HandleAPIError(w, http.StatusForbidden, "API Unavailable")
				return
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			token := ""
			if strings.HasPrefix(authHeader, BearerTokenPrefix) {
				token = strings.TrimPrefix(authHeader, BearerTokenPrefix)
			} else if customKey := r.Header.Get("X-API-Key"); customKey != "" {
				token = customKey
			}

			if token == "" {
				log.Printf("API auth failed: Token missing")
				handlers.HandleAPIError(w, http.StatusUnauthorized, "Authentication Failed")
				return
			}

			hasher := sha256.New()
			hasher.Write([]byte(token))
			incomingHashBytes := hasher.Sum(nil)

			if len(incomingHashBytes) != len(expectedHashedKeyBytes) ||
				subtle.ConstantTimeCompare(incomingHashBytes, expectedHashedKeyBytes) == 0 {
				log.Printf("API auth failed: Invalid token hash mismatch")
				handlers.HandleAPIError(w, http.StatusForbidden, "Authentication Failed")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
