package events

import "time"

const (
	TypePostCreated    = "post_created"
	TypeCommentCreated = "comment_created"
	TypeLikeCreated    = "like_created"
	TypeViewCreated    = "view_created"
)

type Event struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type Data struct {
	PostID    int       `json:"post_id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	Author    string    `json:"author,omitempty"`
	User      string    `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
