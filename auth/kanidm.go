package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/scenery/mediax/config"
	"golang.org/x/oauth2"
)

var (
	oauth2Config *oauth2.Config
	oidcProvider *oidc.Provider
	oidcVerifier *oidc.IDTokenVerifier

	// Store PKCE verifiers temporarily (keyed by state)
	// Using oauth2.GenerateVerifier() to generate verifiers
	pkceVerifiers      = make(map[string]string)
	pkceVerifiersMutex sync.Mutex
)

// InitKanidm initializes the Kanidm OIDC provider and OAuth2 configuration
func InitKanidm() error {
	if !config.App.Kanidm.Enabled {
		log.Println("Kanidm authentication is disabled")
		return nil
	}

	ctx := context.Background()

	// Initialize OIDC provider with auto-discovery
	provider, err := oidc.NewProvider(ctx, config.App.Kanidm.IssuerURL)
	if err != nil {
		return fmt.Errorf("failed to initialize Kanidm OIDC provider: %w", err)
	}
	oidcProvider = provider

	// Create OAuth2 configuration
	oauth2Config = &oauth2.Config{
		ClientID:     config.App.Kanidm.ClientID,
		ClientSecret: config.App.Kanidm.ClientSecret,
		RedirectURL:  config.App.Kanidm.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       config.App.Kanidm.Scopes,
	}

	// Create ID token verifier
	oidcVerifier = provider.Verifier(&oidc.Config{
		ClientID: config.App.Kanidm.ClientID,
	})

	log.Printf("Kanidm OIDC provider initialized: %s", config.App.Kanidm.IssuerURL)
	return nil
}

// GetKanidmAuthURL generates the authorization URL for Kanidm login with PKCE
func GetKanidmAuthURL(state string) string {
	if oauth2Config == nil {
		log.Println("Warning: Kanidm OAuth2 config not initialized")
		return ""
	}

	// Generate PKCE verifier using built-in oauth2 function
	// This generates a fresh verifier with 32 octets of randomness (RFC 7636)
	verifier := oauth2.GenerateVerifier()

	// Store verifier for later use during token exchange (keyed by state)
	pkceVerifiersMutex.Lock()
	pkceVerifiers[state] = verifier
	pkceVerifiersMutex.Unlock()

	// Generate authorization URL with PKCE S256 challenge
	// oauth2.S256ChallengeOption automatically derives the challenge from verifier
	return oauth2Config.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
}

// ExchangeKanidmCode exchanges the authorization code for tokens with PKCE verifier
func ExchangeKanidmCode(ctx context.Context, code string, state string) (*oauth2.Token, error) {
	if oauth2Config == nil {
		return nil, fmt.Errorf("Kanidm OAuth2 config not initialized")
	}

	// Retrieve the PKCE verifier for this state
	pkceVerifiersMutex.Lock()
	verifier, ok := pkceVerifiers[state]
	if ok {
		delete(pkceVerifiers, state) // Clean up after use
	}
	pkceVerifiersMutex.Unlock()

	if !ok {
		return nil, fmt.Errorf("PKCE verifier not found for state")
	}

	// Exchange authorization code for tokens with PKCE verifier
	// oauth2.VerifierOption passes the verifier to the token endpoint
	token, err := oauth2Config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return token, nil
}

// VerifyIDToken verifies and extracts claims from the ID token
func VerifyIDToken(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	if oidcVerifier == nil {
		return nil, fmt.Errorf("Kanidm OIDC verifier not initialized")
	}

	// Verify ID token
	idToken, err := oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	return idToken, nil
}

// GetUserInfoFromIDToken extracts user information from ID token
func GetUserInfoFromIDToken(idToken *oidc.IDToken) (map[string]interface{}, error) {
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to extract claims: %w", err)
	}
	return claims, nil
}

// GenerateStateToken generates a random state token for CSRF protection
func GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// IsKanidmEnabled returns whether Kanidm authentication is enabled
func IsKanidmEnabled() bool {
	return config.App.Kanidm.Enabled && oauth2Config != nil
}
