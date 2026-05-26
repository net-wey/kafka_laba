package transporthttp

import (
	"github.com/gofiber/fiber/v2"
	apisvc "kafka_laba/api-service/service"
)

type ProducerSender = apisvc.ProducerSender
type DataGetter = apisvc.DataGetter

func RegisterHTTP(app *fiber.App, producer ProducerSender, client DataGetter) {
	h := NewWebHandler(apisvc.New(producer, client))

	group := app.Group("/api/v1")
	group.Post("/posts", h.PostCreate)
	group.Get("/posts", h.PostList)
	group.Get("/posts/:id/comments", h.PostComments)
	group.Post("/comments", h.CommentCreate)
	group.Post("/likes", h.PostLike)
	group.Post("/views", h.PostView)
	group.Get("/search", h.Search)
	group.Get("/reports/top-posts-by-comments", h.ReportTopPostsByComments)
	group.Get("/reports/posts-comments-by-day", h.ReportPostsCommentsByDay)
	group.Get("/reports/top-posts-by-views", h.ReportTopPostsByViews)
}
