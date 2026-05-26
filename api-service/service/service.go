package service

import (
	"context"
	"time"

	events "kafka_laba/api-service/contracts"
)

type ProducerSender interface {
	Send(ctx context.Context, eventType string, payload events.Data) error
}

type DataGetter interface {
	Get(path string, query map[string]string) ([]byte, int, error)
}

type Service struct {
	producer ProducerSender
	client   DataGetter
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func New(producer ProducerSender, client DataGetter) *Service {
	return &Service{
		producer: producer,
		client:   client,
	}
}

func (s *Service) CreatePost(ctx context.Context, title, body, author string) error {
	if title == "" || body == "" || author == "" {
		return &ValidationError{Message: "title, body, author are required"}
	}

	return s.producer.Send(ctx, events.TypePostCreated, events.Data{
		Title:     title,
		Body:      body,
		Author:    author,
		CreatedAt: time.Now().UTC(),
	})
}

func (s *Service) CreateComment(ctx context.Context, postID int, body, author string) error {
	if postID <= 0 || body == "" || author == "" {
		return &ValidationError{Message: "post_id, body, author are required"}
	}

	return s.producer.Send(ctx, events.TypeCommentCreated, events.Data{
		PostID:    postID,
		Body:      body,
		Author:    author,
		CreatedAt: time.Now().UTC(),
	})
}

func (s *Service) CreateLike(ctx context.Context, postID int, user string) error {
	if postID <= 0 || user == "" {
		return &ValidationError{Message: "post_id and user are required"}
	}

	return s.producer.Send(ctx, events.TypeLikeCreated, events.Data{
		PostID:    postID,
		User:      user,
		CreatedAt: time.Now().UTC(),
	})
}

func (s *Service) CreateView(ctx context.Context, postID int, user string) error {
	if postID <= 0 {
		return &ValidationError{Message: "post_id is required"}
	}

	return s.producer.Send(ctx, events.TypeViewCreated, events.Data{
		PostID:    postID,
		User:      user,
		CreatedAt: time.Now().UTC(),
	})
}

func (s *Service) ListPosts(limit string) ([]byte, int, error) {
	return s.client.Get("/internal/v1/posts", map[string]string{"limit": withDefault(limit, "20")})
}

func (s *Service) ListPostComments(postID, limit string) ([]byte, int, error) {
	return s.client.Get("/internal/v1/posts/"+postID+"/comments", map[string]string{"limit": withDefault(limit, "100")})
}

func (s *Service) Search(query string) ([]byte, int, error) {
	return s.client.Get("/internal/v1/search", map[string]string{"query": query})
}

func (s *Service) ReportTopPostsByComments(limit string) ([]byte, int, error) {
	return s.client.Get("/internal/v1/reports/top-posts-by-comments", map[string]string{"limit": withDefault(limit, "10")})
}

func (s *Service) ReportPostsCommentsByDay() ([]byte, int, error) {
	return s.client.Get("/internal/v1/reports/posts-comments-by-day", nil)
}

func (s *Service) ReportTopPostsByViews(limit string) ([]byte, int, error) {
	return s.client.Get("/internal/v1/reports/top-posts-by-views", map[string]string{"limit": withDefault(limit, "10")})
}

func withDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
