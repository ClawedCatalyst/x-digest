package postgres

import (
	"context"
	_ "embed"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

// DB holds the connection pool and runs migrations.
type DB struct {
	Pool *pgxpool.Pool
}

// NewDB creates a connection pool and runs the schema.
func NewDB(ctx context.Context, url string) (*DB, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	d := &DB{Pool: pool}
	if err := d.ExecSchema(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return d, nil
}

// Close closes the pool.
func (d *DB) Close() { d.Pool.Close() }

// ExecSchema runs the embedded schema SQL.
func (d *DB) ExecSchema(ctx context.Context) error {
	_, err := d.Pool.Exec(ctx, schemaSQL)
	return err
}

// UserRepo implements ports.UserRepository.
type UserRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo returns a UserRepo.
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// Upsert inserts or updates a user by x_user_id, returns id.
func (r *UserRepo) Upsert(ctx context.Context, xUserID, username string) (string, error) {
	q := `
insert into app_users (x_user_id, username)
values ($1, $2)
on conflict (x_user_id) do update set username = excluded.username
returning id;
`
	var id string
	err := r.pool.QueryRow(ctx, q, xUserID, username).Scan(&id)
	return id, err
}

// TokenRepo implements ports.TokenRepository.
type TokenRepo struct {
	pool *pgxpool.Pool
}

// NewTokenRepo returns a TokenRepo.
func NewTokenRepo(pool *pgxpool.Pool) *TokenRepo {
	return &TokenRepo{pool: pool}
}

// Upsert inserts or updates tokens for a user.
func (r *TokenRepo) Upsert(ctx context.Context, userID, accessEnc, refreshEnc string, expiresAt time.Time) error {
	q := `
insert into tokens (user_id, access_token_enc, refresh_token_enc, expires_at)
values ($1, $2, nullif($3,''), $4)
on conflict (user_id) do update set
  access_token_enc = excluded.access_token_enc,
  refresh_token_enc = excluded.refresh_token_enc,
  expires_at = excluded.expires_at,
  updated_at = now();
`
	_, err := r.pool.Exec(ctx, q, userID, accessEnc, refreshEnc, expiresAt)
	return err
}

// GetByUserID returns encrypted tokens and x_user_id for the user.
func (r *TokenRepo) GetByUserID(ctx context.Context, userID string) (accessEnc, refreshEnc string, expiresAt time.Time, xUserID string, err error) {
	q := `
select t.access_token_enc, coalesce(t.refresh_token_enc,''), t.expires_at, u.x_user_id
from tokens t join app_users u on u.id = t.user_id
where t.user_id = $1
`
	err = r.pool.QueryRow(ctx, q, userID).Scan(&accessEnc, &refreshEnc, &expiresAt, &xUserID)
	return accessEnc, refreshEnc, expiresAt, xUserID, err
}
