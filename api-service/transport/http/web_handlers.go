package transporthttp

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	apisvc "kafka_laba/api-service/service"
)

type WebHandler struct {
	service *apisvc.Service
}

func NewWebHandler(service *apisvc.Service) *WebHandler {
	return &WebHandler{service: service}
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

func (h *WebHandler) PostCreate(c *fiber.Ctx) error {
	var req PostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := h.service.CreatePost(c.UserContext(), req.Title, req.Body, req.Author); err != nil {
		return writeCreateError(c, err)
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "queued"})
}

func (h *WebHandler) PostList(c *fiber.Ctx) error {
	body, status, err := h.service.ListPosts(c.Query("limit", "20"))
	if err != nil {
		return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusOK).Send(body)
}

func (h *WebHandler) PostComments(c *fiber.Ctx) error {
	body, status, err := h.service.ListPostComments(c.Params("id"), c.Query("limit", "100"))
	if err != nil {
		return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusOK).Send(body)
}

func (h *WebHandler) CommentCreate(c *fiber.Ctx) error {
	var req CommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := h.service.CreateComment(c.UserContext(), req.PostID, req.Body, req.Author); err != nil {
		return writeCreateError(c, err)
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "queued"})
}

func (h *WebHandler) PostLike(c *fiber.Ctx) error {
	var req LikeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := h.service.CreateLike(c.UserContext(), req.PostID, req.User); err != nil {
		return writeCreateError(c, err)
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "queued"})
}

func (h *WebHandler) PostView(c *fiber.Ctx) error {
	var req ViewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := h.service.CreateView(c.UserContext(), req.PostID, req.User); err != nil {
		return writeCreateError(c, err)
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "queued"})
}

func (h *WebHandler) Search(c *fiber.Ctx) error {
	body, status, err := h.service.Search(c.Query("query"))
	if err != nil {
		return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusOK).Send(body)
}

func (h *WebHandler) ReportTopPostsByComments(c *fiber.Ctx) error {
	body, status, err := h.service.ReportTopPostsByComments(c.Query("limit", "10"))
	if err != nil {
		return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusOK).Send(body)
}

func (h *WebHandler) ReportPostsCommentsByDay(c *fiber.Ctx) error {
	body, status, err := h.service.ReportPostsCommentsByDay()
	if err != nil {
		return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusOK).Send(body)
}

func (h *WebHandler) ReportTopPostsByViews(c *fiber.Ctx) error {
	body, status, err := h.service.ReportTopPostsByViews(c.Query("limit", "10"))
	if err != nil {
		return c.Status(statusOrDefault(status, fiber.StatusBadGateway)).JSON(fiber.Map{"error": err.Error()})
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusOK).Send(body)
}

func statusOrDefault(status, fallback int) int {
	if status == 0 {
		return fallback
	}
	return status
}

func writeCreateError(c *fiber.Ctx, err error) error {
	var validationErr *apisvc.ValidationError
	if errors.As(err, &validationErr) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Message})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}
