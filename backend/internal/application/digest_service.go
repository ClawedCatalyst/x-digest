package application

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"xdigest/internal/domain"
	"xdigest/internal/ports"
)

// DigestService implements the digest use cases.
type DigestService struct {
	digest ports.DigestRepository
	tok    ports.TokenRepository
	crypt  ports.Crypter
	x      ports.XAPI
}

// NewDigestService constructs a DigestService.
func NewDigestService(digest ports.DigestRepository, tok ports.TokenRepository, crypt ports.Crypter, x ports.XAPI) *DigestService {
	return &DigestService{digest: digest, tok: tok, crypt: crypt, x: x}
}

// GetDigest returns the stored digest for the user and day if present.
func (s *DigestService) GetDigest(ctx context.Context, userID string, day time.Time) ([]byte, bool, error) {
	return s.digest.Get(ctx, userID, day)
}

// BuildDigest fetches data from X, computes the digest, and persists it.
func (s *DigestService) BuildDigest(ctx context.Context, userID string, day time.Time) error {
	accessEnc, refreshEnc, exp, xUserID, err := s.tok.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	accessToken, err := s.crypt.Decrypt(accessEnc)
	if err != nil {
		return err
	}
	refreshToken := ""
	if refreshEnc != "" {
		refreshToken, _ = s.crypt.Decrypt(refreshEnc)
	}

	// Refresh if nearly expired
	if time.Now().After(exp.Add(-30*time.Second)) && refreshToken != "" {
		newTok, err := s.x.RefreshToken(ctx, refreshToken)
		if err != nil {
			return err
		}
		accessToken = newTok.AccessToken
		accessEnc, _ = s.crypt.Encrypt(newTok.AccessToken)
		refreshEnc = ""
		if newTok.RefreshToken != "" {
			refreshEnc, _ = s.crypt.Encrypt(newTok.RefreshToken)
		} else {
			refreshEnc, _ = s.crypt.Encrypt(refreshToken)
		}
		exp = time.Now().Add(time.Duration(newTok.ExpiresIn) * time.Second)
		if err := s.tok.Upsert(ctx, userID, accessEnc, refreshEnc, exp); err != nil {
			return err
		}
	}

	start := day
	end := day.Add(24 * time.Hour)

	posts, err := s.x.GetMyPostsToday(ctx, accessToken, xUserID, start, end)
	if err != nil {
		return err
	}
	mentions, err := s.x.GetMentionsToday(ctx, accessToken, xUserID, start, end)
	if err != nil {
		return err
	}

	eng := map[string]*domain.EngagerStat{}
	usernameByID := map[string]string{}

	for _, m := range mentions {
		if m.AuthorID == "" {
			continue
		}
		es := eng[m.AuthorID]
		if es == nil {
			es = &domain.EngagerStat{UserID: m.AuthorID}
			eng[m.AuthorID] = es
		}
		es.Mentions++
	}

	yesterday := day.Add(-24 * time.Hour)
	var newLikes []domain.LikeEvent
	for _, p := range posts {
		likers, err := s.x.GetLikingUsers(ctx, accessToken, p.ID)
		if err != nil {
			return fmt.Errorf("likers for %s: %w", p.ID, err)
		}
		todaySet := map[string]domain.UserLite{}
		for _, u := range likers {
			todaySet[u.ID] = u
			usernameByID[u.ID] = u.Username
		}
		yset, _ := s.digest.GetLikerSnapshot(ctx, userID, yesterday, p.ID)
		var newUsernames []string
		for id, u := range todaySet {
			if _, ok := yset[id]; !ok {
				newUsernames = append(newUsernames, u.Username)
				es := eng[id]
				if es == nil {
					es = &domain.EngagerStat{UserID: id}
					eng[id] = es
				}
				es.Likes++
			}
		}
		if len(newUsernames) > 0 {
			sort.Strings(newUsernames)
			newLikes = append(newLikes, domain.LikeEvent{TweetID: p.ID, Users: newUsernames})
		}
		_ = s.digest.PutLikerSnapshot(ctx, userID, day, p.ID, todaySet)
	}

	var newReplies []domain.ReplyEvent
	for _, p := range posts {
		conv := p.ConversationID
		if conv == "" {
			conv = p.ID
		}
		replies, err := s.x.GetRepliesTodayByConversation(ctx, accessToken, conv, start, end)
		if err != nil {
			return fmt.Errorf("replies for %s: %w", p.ID, err)
		}
		for _, rt := range replies {
			if rt.ID == p.ID {
				continue
			}
			if rt.AuthorID == "" || rt.AuthorID == xUserID {
				continue
			}
			newReplies = append(newReplies, domain.ReplyEvent{
				RootTweetID: p.ID,
				ReplyID:     rt.ID,
				AuthorID:    rt.AuthorID,
				Text:        rt.Text,
				CreatedAt:   rt.CreatedAt.Format(time.RFC3339),
			})
			es := eng[rt.AuthorID]
			if es == nil {
				es = &domain.EngagerStat{UserID: rt.AuthorID}
				eng[rt.AuthorID] = es
			}
			es.Replies++
		}
	}

	for id, es := range eng {
		if es.Username == "" {
			es.Username = usernameByID[id]
		}
		es.Total = es.Likes + es.Replies + es.Mentions
	}

	top := make([]domain.EngagerStat, 0, len(eng))
	for _, v := range eng {
		top = append(top, *v)
	}
	sort.Slice(top, func(i, j int) bool { return top[i].Total > top[j].Total })
	if len(top) > 10 {
		top = top[:10]
	}

	d := domain.Digest{
		Day:         day.Format("2006-01-02"),
		PostsToday:  posts,
		Mentions:    mentions,
		NewLikes:    newLikes,
		NewReplies:  newReplies,
		TopEngagers: top,
	}
	b, _ := json.Marshal(d)
	return s.digest.Save(ctx, userID, day, b)
}
