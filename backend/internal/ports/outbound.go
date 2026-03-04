package ports

import (
	"context"
	"time"

	"xdigest/internal/domain"
)

// Crypter encrypts/decrypts sensitive strings (tokens, session).
type Crypter interface {
	Encrypt(plain string) (string, error)
	Decrypt(cipher string) (string, error)
}

// UserRepository persists app users.
type UserRepository interface {
	Upsert(ctx context.Context, xUserID, username string) (string, error)
}

// TokenRepository persists encrypted OAuth tokens per user.
type TokenRepository interface {
	Upsert(ctx context.Context, userID, accessEnc, refreshEnc string, expiresAt time.Time) error
	GetByUserID(ctx context.Context, userID string) (accessEnc, refreshEnc string, expiresAt time.Time, xUserID string, err error)
}

// DigestRepository persists and retrieves daily digest data and like snapshots.
type DigestRepository interface {
	Get(ctx context.Context, userID string, day time.Time) (data []byte, found bool, err error)
	Save(ctx context.Context, userID string, day time.Time, data []byte) error
	GetLikerSnapshot(ctx context.Context, userID string, day time.Time, tweetID string) (map[string]domain.UserLite, error)
	PutLikerSnapshot(ctx context.Context, userID string, day time.Time, tweetID string, snapshot map[string]domain.UserLite) error
}

// XAPI is the outbound port for the X (Twitter) API.
type XAPI interface {
	ExchangeCode(ctx context.Context, code, codeVerifier string) (*domain.TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenResponse, error)
	GetMe(ctx context.Context, accessToken string) (*MeResponse, error)
	GetMyPostsToday(ctx context.Context, accessToken, userID string, startUTC, endUTC time.Time) ([]domain.Tweet, error)
	GetMentionsToday(ctx context.Context, accessToken, userID string, startUTC, endUTC time.Time) ([]domain.Tweet, error)
	GetLikingUsers(ctx context.Context, accessToken, tweetID string) ([]domain.UserLite, error)
	GetRepliesTodayByConversation(ctx context.Context, accessToken, conversationID string, startUTC, endUTC time.Time) ([]domain.Tweet, error)
}

// MeResponse is the /2/users/me response (ID and username).
type MeResponse struct {
	Data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"data"`
}
