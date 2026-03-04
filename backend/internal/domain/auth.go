package domain

import "time"

// AuthState holds OAuth PKCE state during the callback flow.
type AuthState struct {
	State        string `json:"state"`
	CodeVerifier string `json:"verifier"`
	CreatedAt    int64  `json:"created_at"`
}

// TokenResponse is the OAuth2 token response from X (Twitter).
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}
