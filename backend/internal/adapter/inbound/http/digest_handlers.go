package http

import (
	"net/http"
	"time"
)

// HandleGetDigestToday returns the digest for today for the current user.
func (s *Server) HandleGetDigestToday(w http.ResponseWriter, r *http.Request) {
	userID, err := s.requireUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	day := time.Now().UTC().Truncate(24 * time.Hour)

	data, found, err := s.digest.GetDigest(r.Context(), userID, day)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if !found {
		http.Error(w, "digest not generated yet. call POST /jobs/digest/today", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// HandleBuildDigestToday builds and stores the digest for today.
func (s *Server) HandleBuildDigestToday(w http.ResponseWriter, r *http.Request) {
	userID, err := s.requireUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	day := time.Now().UTC().Truncate(24 * time.Hour)
	if err := s.digest.BuildDigest(r.Context(), userID, day); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}
