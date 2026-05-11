package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"kafka_laba/producer/internal/shared/events"
)

type ProducerSender interface {
	Send(ctx context.Context, eventType string, payload events.Data) error
}

type DataGetter interface {
	Get(path string, query map[string]string) ([]byte, int, error)
}

type PostRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

type CommentRequest struct {
	PostID int    `json:"post_id"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

type LikeRequest struct {
	PostID int    `json:"post_id"`
	User   string `json:"user"`
}

type ViewRequest struct {
	PostID int    `json:"post_id"`
	User   string `json:"user"`
}

func RegisterHTTP(app *fiber.App, producer ProducerSender, client DataGetter) {
	group := app.Group("/api/v1")

	group.Post("/posts", func(c *fiber.Ctx) error {
		var req PostRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		if req.Title == "" || req.Body == "" || req.Author == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title, body, author are required"})
		}

		err := producer.Send(c.UserContext(), events.TypePostCreated, events.Data{
			Title:     req.Title,
			Body:      req.Body,
			Author:    req.Author,
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "queued"})
	})

	group.Get("/posts", func(c *fiber.Ctx) error {
		body, status, err := client.Get("/internal/v1/posts", map[string]string{"limit": c.Query("limit", "20")})
		if err != nil {
			return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Status(fiber.StatusOK).Send(body)
	})

	group.Get("/posts/:id/comments", func(c *fiber.Ctx) error {
		body, status, err := client.Get("/internal/v1/posts/"+c.Params("id")+"/comments", map[string]string{"limit": c.Query("limit", "100")})
		if err != nil {
			return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Status(fiber.StatusOK).Send(body)
	})

	group.Post("/comments", func(c *fiber.Ctx) error {
		var req CommentRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		if req.PostID <= 0 || req.Body == "" || req.Author == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "post_id, body, author are required"})
		}

		err := producer.Send(c.UserContext(), events.TypeCommentCreated, events.Data{
			PostID:    req.PostID,
			Body:      req.Body,
			Author:    req.Author,
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "queued"})
	})

	group.Post("/likes", func(c *fiber.Ctx) error {
		var req LikeRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		if req.PostID <= 0 || req.User == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "post_id and user are required"})
		}

		err := producer.Send(c.UserContext(), events.TypeLikeCreated, events.Data{
			PostID:    req.PostID,
			User:      req.User,
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "queued"})
	})

	group.Post("/views", func(c *fiber.Ctx) error {
		var req ViewRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		if req.PostID <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "post_id is required"})
		}

		err := producer.Send(c.UserContext(), events.TypeViewCreated, events.Data{
			PostID:    req.PostID,
			User:      req.User,
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "queued"})
	})

	group.Get("/search", func(c *fiber.Ctx) error {
		body, status, err := client.Get("/internal/v1/search", map[string]string{"query": c.Query("query")})
		if err != nil {
			return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Status(fiber.StatusOK).Send(body)
	})

	group.Get("/reports/top-posts-by-comments", func(c *fiber.Ctx) error {
		body, status, err := client.Get("/internal/v1/reports/top-posts-by-comments", map[string]string{"limit": c.Query("limit", "10")})
		if err != nil {
			return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Status(fiber.StatusOK).Send(body)
	})

	group.Get("/reports/posts-comments-by-day", func(c *fiber.Ctx) error {
		body, status, err := client.Get("/internal/v1/reports/posts-comments-by-day", nil)
		if err != nil {
			return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Status(fiber.StatusOK).Send(body)
	})

	group.Get("/reports/top-posts-by-views", func(c *fiber.Ctx) error {
		body, status, err := client.Get("/internal/v1/reports/top-posts-by-views", map[string]string{"limit": c.Query("limit", "10")})
		if err != nil {
			return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Status(fiber.StatusOK).Send(body)
	})
}

func statusOrDefault(status, fallback int) int {
	if status == 0 {
		return fallback
	}
	return status
}
