package main

import (
	"log"
	"strings"

	env "kafka_laba/api-service/config"
	events "kafka_laba/api-service/contracts"
	"kafka_laba/api-service/dataclient"
	"kafka_laba/api-service/kafka"
	transporthttp "kafka_laba/api-service/transport/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	port := env.Get("API_PORT", "8080")
	brokersRaw := env.Get("KAFKA_BROKERS", "localhost:9092")
	dataURL := env.Get("DATA_SERVICE_URL", "http://localhost:8081")
	postCreatedTopic := env.Get("KAFKA_TOPIC_POST_CREATED", "post-created-events")
	commentCreatedTopic := env.Get("KAFKA_TOPIC_COMMENT_CREATED", "comment-created-events")
	likeCreatedTopic := env.Get("KAFKA_TOPIC_LIKE_CREATED", "like-created-events")
	viewCreatedTopic := env.Get("KAFKA_TOPIC_VIEW_CREATED", "view-created-events")

	producer := kafka.NewProducer(strings.Split(brokersRaw, ","), map[string]string{
		events.TypePostCreated:    postCreatedTopic,
		events.TypeCommentCreated: commentCreatedTopic,
		events.TypeLikeCreated:    likeCreatedTopic,
		events.TypeViewCreated:    viewCreatedTopic,
	})
	defer producer.Close()

	client := dataclient.NewDataClient(dataURL)

	app := fiber.New()
	transporthttp.RegisterHTTP(app, producer, client)

	log.Printf("api-service listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("api-service listen: %v", err)
	}
}
