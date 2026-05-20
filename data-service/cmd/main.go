package main

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	env "kafka_laba/data-service/config"
	"kafka_laba/data-service/data"
	"kafka_laba/data-service/data/ent"
)

func main() {
	ctx := context.Background()

	port := env.Get("DATA_PORT", "8081")
	dsn := env.Get("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/kafka_laba?sslmode=disable")
	brokersRaw := env.Get("KAFKA_BROKERS", "localhost:9092")
	topicsRaw := env.Get("KAFKA_TOPICS", "post-created-events,comment-created-events,like-created-events,view-created-events")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	driver := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(driver))
	defer client.Close()

	if err = ensureSchema(ctx, db); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}

	svc := data.NewService(client, db)
	consumer := data.NewConsumer(strings.Split(brokersRaw, ","), splitAndTrim(topicsRaw), svc)
	defer consumer.Close()

	go consumer.Run(ctx)

	app := fiber.New()
	data.RegisterHTTP(app, svc)

	log.Printf("data-service listening on :%s", port)
	if err = app.Listen(":" + port); err != nil {
		log.Fatalf("data-service listen: %v", err)
	}
}

func splitAndTrim(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			out = append(out, value)
		}
	}
	return out
}

func ensureSchema(ctx context.Context, db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			body TEXT NOT NULL,
			author TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INTEGER NOT NULL REFERENCES posts(id),
			body TEXT NOT NULL,
			author TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS post_likes (
			id SERIAL PRIMARY KEY,
			post_id INTEGER NOT NULL REFERENCES posts(id),
			"user" TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			UNIQUE(post_id, "user")
		)`,
		`CREATE TABLE IF NOT EXISTS post_views (
			id SERIAL PRIMARY KEY,
			post_id INTEGER NOT NULL REFERENCES posts(id),
			"user" TEXT,
			created_at TIMESTAMPTZ NOT NULL
		)`,
	}
	for _, query := range queries {
		if _, err := db.ExecContext(ctx, query); err != nil {
			return err
		}
	}
	return nil
}
