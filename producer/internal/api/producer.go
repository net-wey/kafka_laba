package api

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"kafka_laba/producer/internal/shared/events"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			RequiredAcks: kafka.RequireOne,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *Producer) Send(ctx context.Context, eventType string, payload events.Data) error {
	msg := events.Event{Type: eventType, Data: payload}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   buildMessageKey(eventType, payload),
		Value: bytes,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func buildMessageKey(eventType string, payload events.Data) []byte {
	if payload.PostID > 0 {
		return []byte(strconv.Itoa(payload.PostID))
	}
	if payload.Author != "" {
		return []byte(payload.Author)
	}
	return []byte(eventType)
}
