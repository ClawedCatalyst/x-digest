package config

import (
	"encoding/base64"
	"errors"
	"os"
)

// Config holds application configuration from environment.
type Config struct {
	Port            string
	AppBaseURL      string
	FrontendBaseURL string

	XClientID     string
	XClientSecret string
	XRedirectURI  string
	XScopes       string

	CookieSecret []byte
	EncKey       []byte

	DatabaseURL string
}

// Load reads config from environment.
func Load() (*Config, error) {
	get := func(k string) string { return os.Getenv(k) }

	c := &Config{
		Port:            firstNonEmpty(get("PORT"), "8080"),
		AppBaseURL:      firstNonEmpty(get("APP_BASE_URL"), "http://localhost:8080"),
		FrontendBaseURL: firstNonEmpty(get("FRONTEND_BASE_URL"), "http://localhost:3000"),

		XClientID:    get("X_CLIENT_ID"),
		XClientSecret: get("X_CLIENT_SECRET"),
		XRedirectURI:  get("X_REDIRECT_URI"),
		XScopes:       firstNonEmpty(get("X_SCOPES"), "tweet.read users.read like.read offline.access"),

		DatabaseURL: get("DATABASE_URL"),
	}

	if c.XClientID == "" || c.XRedirectURI == "" {
		return nil, errors.New("missing X_CLIENT_ID or X_REDIRECT_URI")
	}
	if get("COOKIE_SECRET") == "" {
		return nil, errors.New("missing COOKIE_SECRET")
	}
	c.CookieSecret = []byte(get("COOKIE_SECRET"))

	encB64 := get("ENCRYPTION_KEY_BASE64")
	if encB64 == "" {
		return nil, errors.New("missing ENCRYPTION_KEY_BASE64")
	}
	enc, err := base64.StdEncoding.DecodeString(encB64)
	if err != nil || len(enc) != 32 {
		return nil, errors.New("ENCRYPTION_KEY_BASE64 must decode to 32 bytes")
	}
	c.EncKey = enc

	if c.DatabaseURL == "" {
		return nil, errors.New("missing DATABASE_URL")
	}

	return c, nil
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
