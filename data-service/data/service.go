package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	events "kafka_laba/data-service/contracts"
	"kafka_laba/data-service/data/ent"
	"kafka_laba/data-service/data/ent/comment"
	"kafka_laba/data-service/data/ent/post"
	"kafka_laba/data-service/data/ent/postlike"
	"kafka_laba/data-service/data/ent/postview"
)

type Service struct {
	client *ent.Client
	db     *sql.DB
}

type SearchResult struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	CommentCount int64  `json:"comment_count"`
	LikeCount    int64  `json:"like_count"`
	ViewCount    int64  `json:"view_count"`
}

type RecentPost struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	CommentCount int64  `json:"comment_count"`
	LikeCount    int64  `json:"like_count"`
	ViewCount    int64  `json:"view_count"`
}

type PostComment struct {
	ID        int64  `json:"id"`
	Author    string `json:"author"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type TopPostByComments struct {
	PostID       int64  `json:"post_id"`
	Title        string `json:"title"`
	CommentCount int64  `json:"comment_count"`
}

type DailyStats struct {
	Day           string `json:"day"`
	PostsCount    int64  `json:"posts_count"`
	CommentsCount int64  `json:"comments_count"`
}

type TopPostByViews struct {
	PostID    int64  `json:"post_id"`
	Title     string `json:"title"`
	ViewCount int64  `json:"view_count"`
}

func NewService(client *ent.Client, db *sql.DB) *Service {
	return &Service{client: client, db: db}
}

func (s *Service) ApplyEvent(ctx context.Context, e events.Event) error {
	switch e.Type {
	case events.TypePostCreated:
		_, err := s.client.Post.Create().
			SetTitle(e.Data.Title).
			SetBody(e.Data.Body).
			SetAuthor(e.Data.Author).
			SetCreatedAt(e.Data.CreatedAt).
			Save(ctx)
		return err
	case events.TypeCommentCreated:
		_, err := s.client.Comment.Create().
			SetPostID(e.Data.PostID).
			SetBody(e.Data.Body).
			SetAuthor(e.Data.Author).
			SetCreatedAt(e.Data.CreatedAt).
			Save(ctx)
		return err
	case events.TypeLikeCreated:
		existing, err := s.client.PostLike.Query().
			Where(
				postlike.PostIDEQ(e.Data.PostID),
				postlike.UserEQ(e.Data.User),
			).
			Only(ctx)
		if err == nil {
			return s.client.PostLike.DeleteOneID(existing.ID).Exec(ctx)
		}
		if !ent.IsNotFound(err) {
			return err
		}

		_, err = s.client.PostLike.Create().
			SetPostID(e.Data.PostID).
			SetUser(e.Data.User).
			SetCreatedAt(e.Data.CreatedAt).
			Save(ctx)
		return err
	case events.TypeViewCreated:
		builder := s.client.PostView.Create().
			SetPostID(e.Data.PostID).
			SetCreatedAt(e.Data.CreatedAt)
		if e.Data.User != "" {
			builder.SetUser(e.Data.User)
		}
		_, err := builder.Save(ctx)
		return err
	default:
		return fmt.Errorf("unknown event type: %s", e.Type)
	}
}

func (s *Service) SearchPosts(ctx context.Context, q string) ([]SearchResult, error) {
	if strings.TrimSpace(q) == "" {
		return []SearchResult{}, nil
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT p.id, p.title, p.author,
       COALESCE((SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id), 0) AS comment_count,
       COALESCE((SELECT COUNT(*) FROM post_likes l WHERE l.post_id = p.id), 0) AS like_count,
       COALESCE((SELECT COUNT(*) FROM post_views v WHERE v.post_id = p.id), 0) AS view_count
FROM posts p
WHERE p.title ILIKE '%' || $1 || '%'
   OR p.body ILIKE '%' || $1 || '%'
   OR EXISTS (SELECT 1 FROM comments c2 WHERE c2.post_id = p.id AND c2.body ILIKE '%' || $1 || '%')
ORDER BY p.id DESC`, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]SearchResult, 0)
	for rows.Next() {
		var item SearchResult
		if scanErr := rows.Scan(&item.ID, &item.Title, &item.Author, &item.CommentCount, &item.LikeCount, &item.ViewCount); scanErr != nil {
			return nil, scanErr
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Service) ListRecentPosts(ctx context.Context, limit int) ([]RecentPost, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `
SELECT p.id, p.title, p.author,
       COALESCE((SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id), 0) AS comment_count,
       COALESCE((SELECT COUNT(*) FROM post_likes l WHERE l.post_id = p.id), 0) AS like_count,
       COALESCE((SELECT COUNT(*) FROM post_views v WHERE v.post_id = p.id), 0) AS view_count
FROM posts p
ORDER BY p.id DESC
LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]RecentPost, 0)
	for rows.Next() {
		var item RecentPost
		if scanErr := rows.Scan(&item.ID, &item.Title, &item.Author, &item.CommentCount, &item.LikeCount, &item.ViewCount); scanErr != nil {
			return nil, scanErr
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Service) ListPostComments(ctx context.Context, postID, limit int) ([]PostComment, error) {
	if postID <= 0 {
		return []PostComment{}, nil
	}
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx, `
SELECT c.id, c.author, c.body, c.created_at
FROM comments c
WHERE c.post_id = $1
ORDER BY c.id DESC
LIMIT $2`, postID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]PostComment, 0)
	for rows.Next() {
		var item PostComment
		if scanErr := rows.Scan(&item.ID, &item.Author, &item.Body, &item.CreatedAt); scanErr != nil {
			return nil, scanErr
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Service) ReportTopPostsByComments(ctx context.Context, limit int) ([]TopPostByComments, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := s.db.QueryContext(ctx, `
SELECT p.id, p.title, COUNT(c.id) AS comment_count
FROM posts p
LEFT JOIN comments c ON c.post_id = p.id
GROUP BY p.id, p.title
ORDER BY comment_count DESC, p.id DESC
LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]TopPostByComments, 0)
	for rows.Next() {
		var item TopPostByComments
		if scanErr := rows.Scan(&item.PostID, &item.Title, &item.CommentCount); scanErr != nil {
			return nil, scanErr
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Service) ReportPostsAndCommentsByDay(ctx context.Context) ([]DailyStats, error) {
	rows, err := s.db.QueryContext(ctx, `
WITH post_stats AS (
	SELECT DATE(created_at) AS day, COUNT(*) AS posts_count
	FROM posts
	GROUP BY DATE(created_at)
), comment_stats AS (
	SELECT DATE(created_at) AS day, COUNT(*) AS comments_count
	FROM comments
	GROUP BY DATE(created_at)
)
SELECT COALESCE(p.day, c.day) AS day,
       COALESCE(p.posts_count, 0) AS posts_count,
       COALESCE(c.comments_count, 0) AS comments_count
FROM post_stats p
FULL OUTER JOIN comment_stats c ON p.day = c.day
ORDER BY day DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]DailyStats, 0)
	for rows.Next() {
		var item DailyStats
		if scanErr := rows.Scan(&item.Day, &item.PostsCount, &item.CommentsCount); scanErr != nil {
			return nil, scanErr
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Service) ReportTopPostsByViews(ctx context.Context, limit int) ([]TopPostByViews, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := s.db.QueryContext(ctx, `
SELECT p.id, p.title, COUNT(v.id) AS view_count
FROM posts p
LEFT JOIN post_views v ON v.post_id = p.id
GROUP BY p.id, p.title
ORDER BY view_count DESC, p.id DESC
LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]TopPostByViews, 0)
	for rows.Next() {
		var item TopPostByViews
		if scanErr := rows.Scan(&item.PostID, &item.Title, &item.ViewCount); scanErr != nil {
			return nil, scanErr
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func CountPosts(ctx context.Context, client *ent.Client) (int, error) {
	return client.Post.Query().Where(post.IDGT(0)).Count(ctx)
}

func CountComments(ctx context.Context, client *ent.Client) (int, error) {
	return client.Comment.Query().Where(comment.IDGT(0)).Count(ctx)
}

func CountLikes(ctx context.Context, client *ent.Client) (int, error) {
	return client.PostLike.Query().Where(postlike.IDGT(0)).Count(ctx)
}

func CountViews(ctx context.Context, client *ent.Client) (int, error) {
	return client.PostView.Query().Where(postview.IDGT(0)).Count(ctx)
}
