SHELL := /bin/bash

GOCACHE ?= $(CURDIR)/.cache/go-build
GOMODCACHE ?= $(CURDIR)/.cache/gomod
GO := GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go

.PHONY: help tidy ent-generate test build up down restart ps logs smoke cli cli-help

help:
	@echo "Available targets:"
	@echo "  make tidy          - go mod tidy in producer and consumer"
	@echo "  make ent-generate  - generate ent code in consumer"
	@echo "  make test          - run tests in producer and consumer"
	@echo "  make build         - build producer and consumer"
	@echo "  make up            - docker compose up --build -d"
	@echo "  make down          - docker compose down"
	@echo "  make restart       - restart stack"
	@echo "  make ps            - show containers"
	@echo "  make logs          - follow logs"
	@echo "  make smoke         - send test data + read reports"
	@echo "  make cli-help      - print CLI usage"
	@echo "  make cli           - run CLI, use ARGS='<command flags>'"

tidy:
	cd producer && $(GO) mod tidy
	cd consumer && $(GO) mod tidy

ent-generate:
	cd consumer && $(GO) generate ./internal/data/ent

test:
	cd producer && $(GO) test ./...
	cd consumer && $(GO) test ./...

build:
	cd producer && $(GO) build ./...
	cd consumer && $(GO) build ./...

up:
	docker compose up --build -d

down:
	docker compose down

restart: down up

ps:
	docker compose ps -a

logs:
	docker compose logs -f --tail=100

smoke:
	curl -s -X POST http://localhost:8080/api/v1/posts -H 'Content-Type: application/json' -d '{"title":"Kafka Post","body":"Hello Kafka","author":"alice"}'
	curl -s -X POST http://localhost:8080/api/v1/comments -H 'Content-Type: application/json' -d '{"post_id":1,"body":"First comment","author":"bob"}'
	curl -s -X POST http://localhost:8080/api/v1/likes -H 'Content-Type: application/json' -d '{"post_id":1,"user":"bob"}'
	curl -s -X POST http://localhost:8080/api/v1/views -H 'Content-Type: application/json' -d '{"post_id":1,"user":"charlie"}'
	sleep 2
	@echo ""
	@echo "SEARCH:"
	curl -s 'http://localhost:8080/api/v1/search?query=Kafka'
	@echo ""
	@echo "REPORT1:"
	curl -s 'http://localhost:8080/api/v1/reports/top-posts-by-comments?limit=10'
	@echo ""
	@echo "REPORT2:"
	curl -s 'http://localhost:8080/api/v1/reports/posts-comments-by-day'
	@echo ""
	@echo "REPORT3:"
	curl -s 'http://localhost:8080/api/v1/reports/top-posts-by-views?limit=10'
	@echo ""

cli-help:
	cd producer && $(GO) run ./cmd/cli help

cli:
	cd producer && $(GO) run ./cmd/cli $(ARGS)
