package data

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func RegisterHTTP(app *fiber.App, service *Service) {
	group := app.Group("/internal/v1")

	group.Get("/posts", func(c *fiber.Ctx) error {
		limit, _ := strconv.Atoi(c.Query("limit", "20"))
		result, err := service.ListRecentPosts(c.UserContext(), limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	group.Get("/posts/:id/comments", func(c *fiber.Ctx) error {
		postID, _ := c.ParamsInt("id")
		limit, _ := strconv.Atoi(c.Query("limit", "100"))
		result, err := service.ListPostComments(c.UserContext(), postID, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	group.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("query")
		result, err := service.SearchPosts(c.UserContext(), query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	group.Get("/reports/top-posts-by-comments", func(c *fiber.Ctx) error {
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		result, err := service.ReportTopPostsByComments(c.UserContext(), limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	group.Get("/reports/posts-comments-by-day", func(c *fiber.Ctx) error {
		result, err := service.ReportPostsAndCommentsByDay(c.UserContext())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	group.Get("/reports/top-posts-by-views", func(c *fiber.Ctx) error {
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		result, err := service.ReportTopPostsByViews(c.UserContext(), limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})
}
