package http

import (
	"net/http"
	"net/url"
	"time"

	"xdigest/internal/application"
	"xdigest/internal/domain"
)

// HandleAuthStart starts the X OAuth flow and redirects to X.
func (s *Server) HandleAuthStart(w http.ResponseWriter, r *http.Request) {
	state := &domain.AuthState{
		State:        application.RandURLSafe(32),
		CodeVerifier: application.RandURLSafe(64),
		CreatedAt:    time.Now().Unix(),
	}
	authURL, cookieValue, err := s.auth.StartAuth(state)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "x_oauth",
		Value:    url.QueryEscape(cookieValue),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, authURL, http.StatusFound)
}

// HandleAuthCallback handles the OAuth callback, exchanges code, and redirects to dashboard.
func (s *Server) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		http.Error(w, "missing code/state", 400)
		return
	}

	c, err := r.Cookie("x_oauth")
	if err != nil {
		http.Error(w, "missing oauth cookie", 400)
		return
	}
	encryptedCookie, _ := url.QueryUnescape(c.Value)

	sessionCookie, err := s.auth.HandleCallback(r.Context(), code, state, encryptedCookie)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    url.QueryEscape(sessionCookie),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "x_oauth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, s.frontendBaseURL+"/dashboard", http.StatusFound)
}
