package data

import "github.com/gofiber/fiber/v2"

func RegisterHTTP(app *fiber.App, service *Service) {
	h := NewWebHandler(service)

	group := app.Group("/internal/v1")
	group.Get("/posts", h.PostList)
	group.Get("/posts/:id/comments", h.PostComments)
	group.Get("/search", h.Search)
	group.Get("/reports/top-posts-by-comments", h.ReportTopPostsByComments)
	group.Get("/reports/posts-comments-by-day", h.ReportPostsCommentsByDay)
	group.Get("/reports/top-posts-by-views", h.ReportTopPostsByViews)
}
