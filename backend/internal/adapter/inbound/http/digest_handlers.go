package http

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"
	"fmt"

	"xdigest/internal/domain"
)

// HandleGetDigestToday returns the digest for today for the current user.
// Kept for backwards compatibility; equivalent to period=daily.
func (s *Server) HandleGetDigestToday(w http.ResponseWriter, r *http.Request) {
	userID, err := s.requireUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	day := time.Now().UTC().Truncate(24 * time.Hour)

	data, found, err := s.digest.GetDigest(r.Context(), userID, day)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "digest not generated yet. call POST /jobs/digest/today", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// HandleGetDigestPeriod returns an aggregated digest for a period:
// period=daily (default), weekly, or monthly.
func (s *Server) HandleGetDigestPeriod(w http.ResponseWriter, r *http.Request) {
	userID, err := s.requireUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "daily"
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	var start, end time.Time

	switch period {
	case "daily":
		start, end = today, today
	case "weekly":
		end = today
		start = today.AddDate(0, 0, -6) // last 7 days including today
	case "monthly":
		end = today
		start = today.AddDate(0, 0, -29) // last 30 days including today
	case "quarterly":
		end = today
		start = today.AddDate(0, 0, -89) // last 90 days including today
	default:
		http.Error(w, "invalid period (use daily, weekly, monthly, or quarterly)", http.StatusBadRequest)
		return
	}

	data, found, err := s.aggregateDigest(r.Context(), userID, start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "digest not generated yet. call POST /jobs/digest/today", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// aggregateDigest loads daily digests in [start, end] and aggregates them into
// a single Digest JSON payload.
func (s *Server) aggregateDigest(ctx context.Context, userID string, start, end time.Time) ([]byte, bool, error) {
	// Collect all daily digests that exist in the range.
	var days []domain.Digest
	for d := start; !d.After(end); d = d.Add(24 * time.Hour) {
		data, found, err := s.digest.GetDigest(ctx, userID, d)
		if err != nil {
			return nil, false, err
		}
		if !found {
			continue
		}
		var dg domain.Digest
		if err := json.Unmarshal(data, &dg); err != nil {
			return nil, false, err
		}
		days = append(days, dg)
	}

	if len(days) == 0 {
		return nil, false, nil
	}

	// If there's only one day, just return it as-is.
	if len(days) == 1 {
		out, err := json.Marshal(days[0])
		if err != nil {
			return nil, false, err
		}
		return out, true, nil
	}

	agg := domain.Digest{
		Day:         start.Format("2006-01-02") + " - " + end.Format("2006-01-02"),
		PeriodStart: start.Format("2006-01-02"),
		PeriodEnd:   end.Format("2006-01-02"),
	}

	// Concatenate per-day slices.
	for _, d := range days {
		agg.PostsToday = append(agg.PostsToday, d.PostsToday...)
		agg.Mentions = append(agg.Mentions, d.Mentions...)
		agg.NewLikes = append(agg.NewLikes, d.NewLikes...)
		agg.NewReplies = append(agg.NewReplies, d.NewReplies...)
	}

	// Aggregate top engagers across days by user.
	byUser := map[string]*domain.EngagerStat{}
	for _, d := range days {
		for _, e := range d.TopEngagers {
			es, ok := byUser[e.UserID]
			if !ok {
				copy := e
				byUser[e.UserID] = &copy
				continue
			}
			es.Likes += e.Likes
			es.Replies += e.Replies
			es.Mentions += e.Mentions
		}
	}

	for _, es := range byUser {
		es.Total = es.Likes + es.Replies + es.Mentions
	}

	agg.TopEngagers = make([]domain.EngagerStat, 0, len(byUser))
	for _, es := range byUser {
		agg.TopEngagers = append(agg.TopEngagers, *es)
	}

	sort.Slice(agg.TopEngagers, func(i, j int) bool {
		return agg.TopEngagers[i].Total > agg.TopEngagers[j].Total
	})
	if len(agg.TopEngagers) > 10 {
		agg.TopEngagers = agg.TopEngagers[:10]
	}

	out, err := json.Marshal(agg)
	if err != nil {
		return nil, false, err
	}
	return out, true, nil
}

// HandleBuildDigestToday builds and stores the digest for today.
// Kept for backwards compatibility; use POST /jobs/digest?period=daily instead.
func (s *Server) HandleBuildDigestToday(w http.ResponseWriter, r *http.Request) {
	userID, err := s.requireUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	day := time.Now().UTC().Truncate(24 * time.Hour)
	if err := s.digest.BuildDigest(r.Context(), userID, day); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// HandleBuildDigestPeriod builds and stores daily digests for the given period.
// POST /jobs/digest?period=daily|weekly|monthly|quarterly
// daily = today only; weekly = last 7 days; monthly = last 30 days; quarterly = last 90 days.
func (s *Server) HandleBuildDigestPeriod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, err := s.requireUser(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "daily"
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	var start, end time.Time

	switch period {
	case "daily":
		start, end = today, today
	case "weekly":
		end = today
		start = today.AddDate(0, 0, -6)
	case "monthly":
		end = today
		start = today.AddDate(0, 0, -29)
	case "quarterly":
		end = today
		start = today.AddDate(0, 0, -89)
	default:
		http.Error(w, "invalid period (use daily, weekly, monthly, or quarterly)", http.StatusBadRequest)
		return
	}

	if err := s.digest.BackfillRange(r.Context(), userID, start, end); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
