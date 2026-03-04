package ports

import (
	"context"
	"time"

	"xdigest/internal/domain"
)

// AuthService is the inbound port for authentication use cases.
type AuthService interface {
	StartAuth(state *domain.AuthState) (authURL string, cookieValue string, err error)
	// HandleCallback exchanges the code for tokens, upserts user and tokens,
	// and returns the encrypted session cookie value and user ID.
	HandleCallback(ctx context.Context, code, state, encryptedCookie string) (sessionCookie string, userID string, err error)
	RequireUser(sessionCookie string) (userID string, err error)
}

// DigestService is the inbound port for digest use cases.
type DigestService interface {
	GetDigest(ctx context.Context, userID string, day time.Time) (data []byte, found bool, err error)
	BuildDigest(ctx context.Context, userID string, day time.Time) error
	// BackfillRange builds digests for all days in [start, end].
	BackfillRange(ctx context.Context, userID string, start, end time.Time) error
}
