package data

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type WebHandler struct {
	service *Service
}

func NewWebHandler(service *Service) *WebHandler {
	return &WebHandler{service: service}
}

func (h *WebHandler) PostList(c *fiber.Ctx) error {
	limit := intQueryOrDefault(c.Query("limit"), 20)
	result, err := h.service.ListRecentPosts(c.UserContext(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *WebHandler) PostComments(c *fiber.Ctx) error {
	postID, _ := c.ParamsInt("id")
	limit := intQueryOrDefault(c.Query("limit"), 100)
	result, err := h.service.ListPostComments(c.UserContext(), postID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *WebHandler) Search(c *fiber.Ctx) error {
	result, err := h.service.SearchPosts(c.UserContext(), c.Query("query"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *WebHandler) ReportTopPostsByComments(c *fiber.Ctx) error {
	limit := intQueryOrDefault(c.Query("limit"), 10)
	result, err := h.service.ReportTopPostsByComments(c.UserContext(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *WebHandler) ReportPostsCommentsByDay(c *fiber.Ctx) error {
	result, err := h.service.ReportPostsAndCommentsByDay(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *WebHandler) ReportTopPostsByViews(c *fiber.Ctx) error {
	limit := intQueryOrDefault(c.Query("limit"), 10)
	result, err := h.service.ReportTopPostsByViews(c.UserContext(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func intQueryOrDefault(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
