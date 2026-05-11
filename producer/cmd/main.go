package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"kafka_laba/producer/internal/api"
	"kafka_laba/producer/internal/shared/env"
)

func main() {
	port := env.Get("API_PORT", "8080")
	brokersRaw := env.Get("KAFKA_BROKERS", "localhost:9092")
	topic := env.Get("KAFKA_TOPIC", "blog-events")
	dataURL := env.Get("DATA_SERVICE_URL", "http://localhost:8081")

	producer := api.NewProducer(strings.Split(brokersRaw, ","), topic)
	defer producer.Close()

	client := api.NewDataClient(dataURL)

	app := fiber.New()
	api.RegisterHTTP(app, producer, client)

	log.Printf("producer listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("producer listen: %v", err)
	}
}
