package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"xdigest/internal/domain"
)

// DigestRepo implements ports.DigestRepository.
type DigestRepo struct {
	pool *pgxpool.Pool
}

// NewDigestRepo returns a DigestRepo.
func NewDigestRepo(pool *pgxpool.Pool) *DigestRepo {
	return &DigestRepo{pool: pool}
}

// Get returns the digest data for the user and day if present.
func (r *DigestRepo) Get(ctx context.Context, userID string, day time.Time) ([]byte, bool, error) {
	q := `select data from daily_digest where user_id = $1 and day = $2`
	var data []byte
	err := r.pool.QueryRow(ctx, q, userID, day).Scan(&data)
	if err != nil {
		return nil, false, nil
	}
	return data, true, nil
}

// Save upserts the digest JSON for the user and day (overwrites data).
func (r *DigestRepo) Save(ctx context.Context, userID string, day time.Time, data []byte) error {
	q := `
insert into daily_digest (user_id, day, data)
values ($1, $2, $3)
on conflict (user_id, day) do update set data = excluded.data;
`
	_, err := r.pool.Exec(ctx, q, userID, day, data)
	return err
}

// PutLikerSnapshot merges a like snapshot for a tweet into the row's data.like_snapshots.
func (r *DigestRepo) PutLikerSnapshot(ctx context.Context, userID string, day time.Time, tweetID string, snapshot map[string]domain.UserLite) error {
	type snap struct {
		IDs map[string]domain.UserLite `json:"ids"`
	}
	b, _ := json.Marshal(snap{IDs: snapshot})
	q := `
insert into daily_digest (user_id, day, data)
values ($1, $2, jsonb_build_object('like_snapshots', jsonb_build_object($3, $4::jsonb)))
on conflict (user_id, day) do update
set data = daily_digest.data || jsonb_build_object('like_snapshots',
  coalesce(daily_digest.data->'like_snapshots','{}'::jsonb) || jsonb_build_object($3, $4::jsonb)
);
`
	_, err := r.pool.Exec(ctx, q, userID, day, tweetID, b)
	return err
}

// GetLikerSnapshot returns the stored like snapshot for the tweet on the given day.
func (r *DigestRepo) GetLikerSnapshot(ctx context.Context, userID string, day time.Time, tweetID string) (map[string]domain.UserLite, error) {
	q := `select data->'like_snapshots'->>$2 from daily_digest where user_id = $1 and day = $3`
	var raw *string
	_ = r.pool.QueryRow(ctx, q, userID, tweetID, day).Scan(&raw)
	if raw == nil || *raw == "" {
		return map[string]domain.UserLite{}, nil
	}
	var parsed struct {
		IDs map[string]domain.UserLite `json:"ids"`
	}
	_ = json.Unmarshal([]byte(*raw), &parsed)
	if parsed.IDs == nil {
		parsed.IDs = map[string]domain.UserLite{}
	}
	return parsed.IDs, nil
}
