package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"xdigest/internal/domain"
	"xdigest/internal/ports"
)

// AuthServiceConfig is the configuration required for AuthService.
type AuthServiceConfig struct {
	XClientID        string
	XRedirectURI     string
	XScopes          string
	FrontendBaseURL  string
}

// AuthService implements the auth use cases.
type AuthService struct {
	cfg   AuthServiceConfig
	crypt ports.Crypter
	users ports.UserRepository
	tok   ports.TokenRepository
	x     ports.XAPI
}

// NewAuthService constructs an AuthService.
func NewAuthService(cfg AuthServiceConfig, crypt ports.Crypter, users ports.UserRepository, tok ports.TokenRepository, x ports.XAPI) *AuthService {
	return &AuthService{cfg: cfg, crypt: crypt, users: users, tok: tok, x: x}
}

// StartAuth starts the OAuth flow: generates state/verifier, returns auth URL and encrypted cookie value.
func (s *AuthService) StartAuth(state *domain.AuthState) (authURL string, cookieValue string, err error) {
	raw, _ := json.Marshal(state)
	enc, err := s.crypt.Encrypt(string(raw))
	if err != nil {
		return "", "", err
	}

	authURLParsed, _ := url.Parse("https://x.com/i/oauth2/authorize")
	q := authURLParsed.Query()
	q.Set("response_type", "code")
	q.Set("client_id", s.cfg.XClientID)
	q.Set("redirect_uri", s.cfg.XRedirectURI)
	q.Set("scope", s.cfg.XScopes)
	q.Set("state", state.State)
	q.Set("code_challenge", pkceChallenge(state.CodeVerifier))
	q.Set("code_challenge_method", "S256")
	authURLParsed.RawQuery = q.Encode()

	return authURLParsed.String(), enc, nil
}

// HandleCallback exchanges the code for tokens, upserts user and tokens, returns session cookie value and user ID.
func (s *AuthService) HandleCallback(ctx context.Context, code, state, encryptedCookie string) (sessionCookie string, userID string, err error) {
	raw, err := s.crypt.Decrypt(encryptedCookie)
	if err != nil {
		return "", "", fmt.Errorf("bad oauth cookie")
	}
	var as domain.AuthState
	if err := json.Unmarshal([]byte(raw), &as); err != nil {
		return "", "", fmt.Errorf("bad oauth cookie")
	}
	if as.State != state {
		return "", "", fmt.Errorf("state mismatch")
	}

	tok, err := s.x.ExchangeCode(ctx, code, as.CodeVerifier)
	if err != nil {
		return "", "", err
	}

	me, err := s.x.GetMe(ctx, tok.AccessToken)
	if err != nil {
		return "", "", err
	}

	userID, err = s.users.Upsert(ctx, me.Data.ID, me.Data.Username)
	if err != nil {
		return "", "", err
	}

	accessEnc, _ := s.crypt.Encrypt(tok.AccessToken)
	refreshEnc := ""
	if tok.RefreshToken != "" {
		refreshEnc, _ = s.crypt.Encrypt(tok.RefreshToken)
	}
	expiresAt := time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	if err := s.tok.Upsert(ctx, userID, accessEnc, refreshEnc, expiresAt); err != nil {
		return "", "", err
	}

	sessionCookie, err = s.crypt.Encrypt(userID)
	if err != nil {
		return "", "", err
	}
	return sessionCookie, userID, nil
}

// RequireUser decrypts the session cookie and returns the user ID.
func (s *AuthService) RequireUser(sessionCookie string) (userID string, err error) {
	if sessionCookie == "" {
		return "", fmt.Errorf("not logged in")
	}
	userID, err = s.crypt.Decrypt(sessionCookie)
	if err != nil {
		return "", fmt.Errorf("invalid session")
	}
	return strings.TrimSpace(userID), nil
}
