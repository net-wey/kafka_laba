# Kafka lab: producer + consumer (Go + Fiber + Ent)

## Запуск

```bash
make up
```

Producer API: `http://localhost:8080`
CLI: `go run ./producer/cmd/cli ...`

Все параметры docker compose вынесены в [.env](/Users/margarita/Desktop/kafka_laba/.env).

## Команды Makefile

```bash
make help
make tidy
make ent-generate
make test
make build
make up
make down
make restart
make ps
make logs
make smoke
make cli-help
make cli ARGS="posts -limit=10"
```

## Docker/env

- Dockerfile producer: [producer/Dockerfile](/Users/margarita/Desktop/kafka_laba/producer/Dockerfile)
- Dockerfile consumer: [consumer/Dockerfile](/Users/margarita/Desktop/kafka_laba/consumer/Dockerfile)
- Изолированный код producer: [producer](/Users/margarita/Desktop/kafka_laba/producer)
- Изолированный код consumer: [consumer](/Users/margarita/Desktop/kafka_laba/consumer)
- Compose: [docker-compose.yml](/Users/margarita/Desktop/kafka_laba/docker-compose.yml)
- Переменные: [.env](/Users/margarita/Desktop/kafka_laba/.env)

## Как это работает

1. Клиент отправляет HTTP-запрос в `producer` (`/api/v1/posts`, `/comments`, `/likes`, `/views`).
2. `producer` не пишет в базу напрямую, а отправляет событие в Kafka topic `blog-events`.
3. `consumer` читает события из Kafka как consumer group `consumer-group`.
4. `consumer` сохраняет данные в PostgreSQL таблицы `posts`, `comments`, `post_likes`, `post_views`.
5. Для поиска и отчетов `producer` проксирует запросы в `consumer`.
6. CLI-клиент (`producer/cmd/cli`) вызывает те же API-эндпоинты и печатает JSON в консоль.
7. У Kafka-сообщений есть `key`: для событий с `post_id` ключ = `post_id`, иначе ключ = `author`.

## CLI (основной способ работы)

```bash
go run ./producer/cmd/cli help
```

Примеры:

```bash
go run ./producer/cmd/cli post -title="Post 1" -body="Hello" -author="alice"
go run ./producer/cmd/cli posts -limit=10
go run ./producer/cmd/cli view -post-id=1 -user="bob"
go run ./producer/cmd/cli like -post-id=1 -user="bob"
go run ./producer/cmd/cli comment -post-id=1 -body="Nice post" -author="bob"
go run ./producer/cmd/cli comments -post-id=1 -limit=20
go run ./producer/cmd/cli search -query="Hello"
go run ./producer/cmd/cli report-comments -limit=10
go run ./producer/cmd/cli report-days
go run ./producer/cmd/cli report-views -limit=10
```

## Producer endpoints

### Добавление порций данных (в Kafka)

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H 'Content-Type: application/json' \
  -d '{"title":"Post 1","body":"Hello","author":"alice"}'

curl -X POST http://localhost:8080/api/v1/comments \
  -H 'Content-Type: application/json' \
  -d '{"post_id":1,"body":"Nice post","author":"bob"}'

curl -X POST http://localhost:8080/api/v1/likes \
  -H 'Content-Type: application/json' \
  -d '{"post_id":1,"user":"bob"}'

curl -X POST http://localhost:8080/api/v1/views \
  -H 'Content-Type: application/json' \
  -d '{"post_id":1,"user":"charlie"}'
```

### Поиск и отчеты (через Consumer)

```bash
curl 'http://localhost:8080/api/v1/search?query=Hello'
curl 'http://localhost:8080/api/v1/reports/top-posts-by-comments?limit=10'
curl 'http://localhost:8080/api/v1/reports/posts-comments-by-day'
curl 'http://localhost:8080/api/v1/reports/top-posts-by-views?limit=10'
```

## Что реализовано по ТЗ

- 5 контейнеров: `producer`, `consumer`, `postgres`, `zookeeper`, `kafka`
- `producer` принимает HTTP и пишет в Kafka
- `consumer` читает Kafka и пишет в PostgreSQL
- минимум 2 связанные таблицы (есть `posts -> comments`, также `posts -> likes/views`)
- 3 отчета:
  - топ постов по количеству комментариев
  - количество постов и комментариев по дням
  - топ постов по просмотрам
# kafka_laba
