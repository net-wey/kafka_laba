package data

import (
	"context"
	"encoding/json"
	events "kafka_laba/data-service/contracts"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	svc    *Service
}

func NewConsumer(brokers []string, topic string, svc *Service) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     "data-service-group",
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})

	return &Consumer{reader: reader, svc: svc}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("kafka read error: %v", err)
			continue
		}

		var evt events.Event
		if err = json.Unmarshal(msg.Value, &evt); err != nil {
			log.Printf("invalid event payload: %v", err)
			continue
		}

		if err = c.svc.ApplyEvent(ctx, evt); err != nil {
			log.Printf("apply event failed: type=%s error=%v", evt.Type, err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
