package http

import (
	"net/http"
	"net/url"

	"xdigest/internal/ports"
)

// Server is the HTTP adapter (inbound). It delegates to application services.
type Server struct {
	auth              ports.AuthService
	digest            ports.DigestService
	frontendBaseURL   string
	sessionCookieName string
}

// NewServer constructs the HTTP server with the given services and config.
func NewServer(auth ports.AuthService, digest ports.DigestService, frontendBaseURL string) *Server {
	return &Server{
		auth:              auth,
		digest:            digest,
		frontendBaseURL:   frontendBaseURL,
		sessionCookieName: "sid",
	}
}

// requireUser reads the session cookie and returns the user ID or an error.
func (s *Server) requireUser(r *http.Request) (string, error) {
	c, err := r.Cookie(s.sessionCookieName)
	if err != nil {
		return "", err
	}
	v, _ := url.QueryUnescape(c.Value)
	return s.auth.RequireUser(v)
}
