package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	env "kafka_laba/api-service/config"
	"kafka_laba/api-service/dataclient"
	"kafka_laba/api-service/kafka"
	transporthttp "kafka_laba/api-service/transport/http"
)

func main() {
	port := env.Get("API_PORT", "8080")
	brokersRaw := env.Get("KAFKA_BROKERS", "localhost:9092")
	topic := env.Get("KAFKA_TOPIC", "blog-events")
	dataURL := env.Get("DATA_SERVICE_URL", "http://localhost:8081")

	producer := kafka.NewProducer(strings.Split(brokersRaw, ","), topic)
	defer producer.Close()

	client := dataclient.NewDataClient(dataURL)

	app := fiber.New()
	transporthttp.RegisterHTTP(app, producer, client)

	log.Printf("api-service listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("api-service listen: %v", err)
	}
}
