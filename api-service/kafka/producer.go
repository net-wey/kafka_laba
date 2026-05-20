package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	events "kafka_laba/api-service/contracts"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writersByEvent map[string]*kafka.Writer
}

func NewProducer(brokers []string, topicsByEvent map[string]string) *Producer {
	writersByEvent := make(map[string]*kafka.Writer, len(topicsByEvent))
	for eventType, topic := range topicsByEvent {
		writersByEvent[eventType] = &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireOne,
			BatchTimeout: 10 * time.Millisecond,
		}
	}
	return &Producer{
		writersByEvent: writersByEvent,
	}
}

func (p *Producer) Send(ctx context.Context, eventType string, payload events.Data) error {
	writer, ok := p.writersByEvent[eventType]
	if !ok {
		return fmt.Errorf("no kafka topic configured for event type: %s", eventType)
	}

	msg := events.Event{Type: eventType, Data: payload}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return writer.WriteMessages(ctx, kafka.Message{
		Key:   buildMessageKey(eventType, payload),
		Value: bytes,
	})
}

func (p *Producer) Close() error {
	var firstErr error
	closedTopics := map[string]struct{}{}
	for _, writer := range p.writersByEvent {
		topic := writer.Topic
		if _, exists := closedTopics[topic]; exists {
			continue
		}
		closedTopics[topic] = struct{}{}
		if err := writer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
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
