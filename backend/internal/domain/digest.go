package domain

import "time"

// Digest is the daily digest payload for a user.
type Digest struct {
	Day         string        `json:"day"`
	PostsToday  []Tweet       `json:"posts_today"`
	Mentions    []Tweet       `json:"mentions_today"`
	NewLikes    []LikeEvent   `json:"new_likes_today"`
	NewReplies  []ReplyEvent  `json:"new_replies_today"`
	TopEngagers []EngagerStat `json:"top_engagers_today"`
}

// LikeEvent represents new likes on a tweet since the previous day.
type LikeEvent struct {
	TweetID string   `json:"tweet_id"`
	Users   []string `json:"usernames"`
}

// ReplyEvent represents a reply to the user's post.
type ReplyEvent struct {
	RootTweetID string `json:"root_tweet_id"`
	ReplyID     string `json:"reply_id"`
	AuthorID    string `json:"author_id"`
	Text        string `json:"text"`
	CreatedAt   string `json:"created_at"`
}

// EngagerStat aggregates engagement (likes, replies, mentions) per user.
type EngagerStat struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Likes    int    `json:"likes"`
	Replies  int    `json:"replies"`
	Mentions int    `json:"mentions"`
	Total    int    `json:"total"`
}

// Tweet is a tweet from the X API.
type Tweet struct {
	ID             string    `json:"id"`
	Text           string    `json:"text"`
	CreatedAt      time.Time `json:"created_at"`
	ConversationID string    `json:"conversation_id,omitempty"`
	AuthorID       string    `json:"author_id,omitempty"`
}

// UserLite is a minimal user from the X API (e.g. liking users).
type UserLite struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
