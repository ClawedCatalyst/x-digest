package xapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"xdigest/internal/domain"
	"xdigest/internal/ports"
)

// Config is the X OAuth and API config.
type Config struct {
	XClientID     string
	XClientSecret string
	XRedirectURI  string
}

// Client implements ports.XAPI.
type Client struct {
	http *http.Client
	cfg  Config
}

// NewClient returns an X API client.
func NewClient(cfg Config) *Client {
	return &Client{
		http: &http.Client{Timeout: 20 * time.Second},
		cfg:  cfg,
	}
}

// Ensure Client implements ports.XAPI.
var _ ports.XAPI = (*Client)(nil)

// ExchangeCode exchanges an auth code for tokens.
func (c *Client) ExchangeCode(ctx context.Context, code, codeVerifier string) (*domain.TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", c.cfg.XRedirectURI)
	form.Set("code_verifier", codeVerifier)
	form.Set("client_id", c.cfg.XClientID)

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.x.com/2/oauth2/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.cfg.XClientSecret != "" {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.cfg.XClientID+":"+c.cfg.XClientSecret)))
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("token exchange failed: %d %s", resp.StatusCode, string(body))
	}
	var tr domain.TokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, err
	}
	return &tr, nil
}

// RefreshToken refreshes the access token.
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	form.Set("client_id", c.cfg.XClientID)

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.x.com/2/oauth2/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.cfg.XClientSecret != "" {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.cfg.XClientID+":"+c.cfg.XClientSecret)))
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("refresh failed: %d %s", resp.StatusCode, string(body))
	}
	var tr domain.TokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, err
	}
	return &tr, nil
}

// GetMe returns the authenticated user.
func (c *Client) GetMe(ctx context.Context, accessToken string) (*ports.MeResponse, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.x.com/2/users/me?user.fields=username,name", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("me failed: %d %s", resp.StatusCode, string(body))
	}
	var out ports.MeResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMyPostsToday returns the user's tweets in the time window.
func (c *Client) GetMyPostsToday(ctx context.Context, accessToken, userID string, startUTC, endUTC time.Time) ([]domain.Tweet, error) {
	u, _ := url.Parse(fmt.Sprintf("https://api.x.com/2/users/%s/tweets", userID))
	q := u.Query()
	q.Set("start_time", startUTC.Format(time.RFC3339))
	q.Set("end_time", endUTC.Format(time.RFC3339))
	q.Set("max_results", "100")
	q.Set("tweet.fields", "created_at,conversation_id,author_id")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("posts failed: %d %s", resp.StatusCode, string(body))
	}
	var out struct {
		Data []domain.Tweet `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetMentionsToday returns mentions in the time window.
func (c *Client) GetMentionsToday(ctx context.Context, accessToken, userID string, startUTC, endUTC time.Time) ([]domain.Tweet, error) {
	u, _ := url.Parse(fmt.Sprintf("https://api.x.com/2/users/%s/mentions", userID))
	q := u.Query()
	q.Set("start_time", startUTC.Format(time.RFC3339))
	q.Set("end_time", endUTC.Format(time.RFC3339))
	q.Set("max_results", "100")
	q.Set("tweet.fields", "created_at,author_id")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("mentions failed: %d %s", resp.StatusCode, string(body))
	}
	var out struct {
		Data []domain.Tweet `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetLikingUsers returns users who liked the tweet.
func (c *Client) GetLikingUsers(ctx context.Context, accessToken, tweetID string) ([]domain.UserLite, error) {
	u := fmt.Sprintf("https://api.x.com/2/tweets/%s/liking_users?user.fields=username,name", tweetID)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("liking_users failed: %d %s", resp.StatusCode, string(body))
	}
	var out struct {
		Data []domain.UserLite `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetRepliesTodayByConversation returns replies in the conversation within the time window.
func (c *Client) GetRepliesTodayByConversation(ctx context.Context, accessToken, conversationID string, startUTC, endUTC time.Time) ([]domain.Tweet, error) {
	u, _ := url.Parse("https://api.x.com/2/tweets/search/recent")
	q := u.Query()
	q.Set("query", "conversation_id:"+conversationID)
	q.Set("max_results", "100")
	q.Set("tweet.fields", "created_at,author_id,conversation_id")
	q.Set("start_time", startUTC.Format(time.RFC3339))
	q.Set("end_time", endUTC.Format(time.RFC3339))
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("search recent failed: %d %s", resp.StatusCode, string(body))
	}
	var out struct {
		Data []domain.Tweet `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
