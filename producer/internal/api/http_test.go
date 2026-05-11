package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"kafka_laba/producer/internal/shared/events"
)

type producerStub struct {
	called bool
	err    error
	typeV  string
	data   events.Data
}

func (p *producerStub) Send(_ context.Context, eventType string, payload events.Data) error {
	p.called = true
	p.typeV = eventType
	p.data = payload
	return p.err
}

type dataClientStub struct {
	body   []byte
	status int
	err    error
	path   string
	query  map[string]string
}

func (d *dataClientStub) Get(path string, query map[string]string) ([]byte, int, error) {
	d.path = path
	d.query = query
	if d.err != nil {
		return nil, d.status, d.err
	}
	return d.body, d.status, nil
}

func TestCreatePostValidation(t *testing.T) {
	app := fiber.New()
	producer := &producerStub{}
	client := &dataClientStub{}
	RegisterHTTP(app, producer, client)

	req := httptest.NewRequest("POST", "/api/v1/posts", bytes.NewBufferString(`{"title":"","body":"x","author":"a"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	if producer.called {
		t.Fatalf("producer should not be called on invalid request")
	}
}

func TestCreatePostAccepted(t *testing.T) {
	app := fiber.New()
	producer := &producerStub{}
	client := &dataClientStub{}
	RegisterHTTP(app, producer, client)

	reqBody := map[string]string{"title": "Hello", "body": "World", "author": "margo"}
	encoded, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusAccepted {
		t.Fatalf("expected 202, got %d", resp.StatusCode)
	}
	if !producer.called {
		t.Fatalf("producer should be called")
	}
	if producer.typeV != events.TypePostCreated {
		t.Fatalf("wrong event type: %s", producer.typeV)
	}
	if producer.data.Title != "Hello" || producer.data.Author != "margo" {
		t.Fatalf("unexpected payload: %+v", producer.data)
	}
}

func TestSearchProxy(t *testing.T) {
	app := fiber.New()
	producer := &producerStub{}
	client := &dataClientStub{body: []byte(`[{"id":1}]`), status: fiber.StatusOK}
	RegisterHTTP(app, producer, client)

	req := httptest.NewRequest("GET", "/api/v1/search?query=hello", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if client.path != "/internal/v1/search" {
		t.Fatalf("wrong path: %s", client.path)
	}
	if client.query["query"] != "hello" {
		t.Fatalf("wrong query: %+v", client.query)
	}
}

func TestListPostsProxy(t *testing.T) {
	app := fiber.New()
	producer := &producerStub{}
	client := &dataClientStub{body: []byte(`[{"id":1}]`), status: fiber.StatusOK}
	RegisterHTTP(app, producer, client)

	req := httptest.NewRequest("GET", "/api/v1/posts?limit=5", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if client.path != "/internal/v1/posts" {
		t.Fatalf("wrong path: %s", client.path)
	}
	if client.query["limit"] != "5" {
		t.Fatalf("wrong limit: %+v", client.query)
	}
}

func TestPostCommentsProxy(t *testing.T) {
	app := fiber.New()
	producer := &producerStub{}
	client := &dataClientStub{body: []byte(`[{"id":11}]`), status: fiber.StatusOK}
	RegisterHTTP(app, producer, client)

	req := httptest.NewRequest("GET", "/api/v1/posts/3/comments?limit=15", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if client.path != "/internal/v1/posts/3/comments" {
		t.Fatalf("wrong path: %s", client.path)
	}
	if client.query["limit"] != "15" {
		t.Fatalf("wrong limit: %+v", client.query)
	}
}

func TestSearchProxyError(t *testing.T) {
	app := fiber.New()
	producer := &producerStub{}
	client := &dataClientStub{status: fiber.StatusBadGateway, err: errors.New("boom")}
	RegisterHTTP(app, producer, client)

	req := httptest.NewRequest("GET", "/api/v1/search?query=hello", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadGateway {
		t.Fatalf("expected 502, got %d", resp.StatusCode)
	}
}
