package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/scenery/mediax/auth"
	"github.com/scenery/mediax/config"
)

const (
	stateSessionKey = "oauth_state"
	stateCookieName = "oauth_state"
)

// handleKanidmLogin initiates the Kanidm OAuth2 login flow
func handleKanidmLogin(w http.ResponseWriter, r *http.Request) {
	if !auth.IsKanidmEnabled() {
		handleError(w, "Kanidm authentication is not enabled", "login", 503)
		return
	}

	// Check if user is already authenticated
	if auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	// Generate CSRF protection state token
	state, err := auth.GenerateStateToken()
	if err != nil {
		log.Printf("Error generating state token: %v", err)
		handleError(w, "Internal Server Error: Failed to initiate login", "login", 500)
		return
	}

	// Store state in cookie for verification in callback
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https",
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to Kanidm authorization URL
	authURL := auth.GetKanidmAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// handleKanidmCallback handles the OAuth2 callback from Kanidm
func handleKanidmCallback(w http.ResponseWriter, r *http.Request) {
	if !auth.IsKanidmEnabled() {
		handleError(w, "Kanidm authentication is not enabled", "login", 503)
		return
	}

	// Verify state for CSRF protection
	stateCookie, err := r.Cookie(stateCookieName)
	if err != nil {
		log.Printf("Error reading state cookie: %v", err)
		handleError(w, "Invalid authentication state", "login", 400)
		return
	}

	stateParam := r.URL.Query().Get("state")
	if stateParam == "" || stateParam != stateCookie.Value {
		log.Printf("State mismatch: cookie=%s, param=%s", stateCookie.Value, stateParam)
		handleError(w, "Invalid authentication state", "login", 400)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Check for OAuth2 errors
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		errDesc := r.URL.Query().Get("error_description")
		log.Printf("OAuth2 error: %s - %s", errMsg, errDesc)
		handleError(w, "Authentication failed: "+errMsg, "login", 401)
		return
	}

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("No authorization code received")
		handleError(w, "Invalid authorization response", "login", 400)
		return
	}

	// Exchange code for tokens with PKCE verifier
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oauth2Token, err := auth.ExchangeKanidmCode(ctx, code, stateParam)
	if err != nil {
		log.Printf("Error exchanging code: %v", err)
		handleError(w, "Authentication failed", "login", 500)
		return
	}

	// Extract and verify ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Println("No id_token in OAuth2 token")
		handleError(w, "Invalid token response", "login", 500)
		return
	}

	idToken, err := auth.VerifyIDToken(ctx, rawIDToken)
	if err != nil {
		log.Printf("Error verifying ID token: %v", err)
		handleError(w, "Invalid authentication token", "login", 401)
		return
	}

	// Extract user information from ID token
	claims, err := auth.GetUserInfoFromIDToken(idToken)
	if err != nil {
		log.Printf("Error extracting claims: %v", err)
		handleError(w, "Failed to retrieve user information", "login", 500)
		return
	}

	// Extract user information
	userID, _ := claims["sub"].(string)
	preferredUsername, _ := claims["preferred_username"].(string)
	email, _ := claims["email"].(string)

	log.Printf("Kanidm authentication successful: sub=%s, preferred_username=%s, email=%s",
		userID, preferredUsername, email)

	// Check if the user is authorized
	// Only allow the user configured in config.json to login via Kanidm
	if preferredUsername != config.App.User.Username {
		log.Printf("Authorization failed: Kanidm user '%s' is not authorized (expected: '%s')",
			preferredUsername, config.App.User.Username)
		handleError(w, "You are not authorized to access this application", "login", 403)
		return
	}

	log.Printf("User '%s' authorized successfully", preferredUsername)

	// Create session for the authenticated user
	err = auth.CreateSession(w)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		handleError(w, "Internal Server Error: Failed to create session", "login", 500)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/home", http.StatusFound)
}
