# Kafka lab: api-service + data-service (Go + Fiber + Ent)

## Запуск

```bash
make up
```

API Service: `http://localhost:8080`
CLI: `go run ./api-service/cmd/cli ...`

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

- Dockerfile api-service: [api-service/Dockerfile](/Users/margarita/Desktop/kafka_laba/api-service/Dockerfile)
- Dockerfile data-service: [data-service/Dockerfile](/Users/margarita/Desktop/kafka_laba/data-service/Dockerfile)
- Изолированный код api-service: [api-service](/Users/margarita/Desktop/kafka_laba/api-service)
- Изолированный код data-service: [data-service](/Users/margarita/Desktop/kafka_laba/data-service)
- Compose: [docker-compose.yml](/Users/margarita/Desktop/kafka_laba/docker-compose.yml)
- Переменные: [.env](/Users/margarita/Desktop/kafka_laba/.env)

## Структура

```text
api-service/
  cmd/
    main.go
    cli/main.go
  transport/http/
  kafka/
  dataclient/
  config/
  contracts/

data-service/
  cmd/
    main.go
  data/
  config/
  contracts/
```

## Как это работает

1. Клиент отправляет HTTP-запрос в `api-service` (`/api/v1/posts`, `/comments`, `/likes`, `/views`).
2. `api-service` не пишет в базу напрямую, а отправляет события в отдельные Kafka topic:
   `post-created-events`, `comment-created-events`, `like-created-events`, `view-created-events`.
3. `data-service` читает события из Kafka как consumer group `data-service-group`.
4. `data-service` сохраняет данные в PostgreSQL таблицы `posts`, `comments`, `post_likes`, `post_views`.
5. Для поиска и отчетов `api-service` проксирует запросы в `data-service`.
6. CLI-клиент (`api-service/cmd/cli`) вызывает те же API-эндпоинты и печатает JSON в консоль.
7. У Kafka-сообщений есть `key`: для событий с `post_id` ключ = `post_id`, иначе ключ = `author`.

## CLI (основной способ работы)

```bash
go run ./api-service/cmd/cli help
```

Примеры:

```bash
go run ./api-service/cmd/cli post -title="Post 1" -body="Hello" -author="alice"
go run ./api-service/cmd/cli posts -limit=10
go run ./api-service/cmd/cli view -post-id=1 -user="bob"
go run ./api-service/cmd/cli like -post-id=1 -user="bob"
go run ./api-service/cmd/cli comment -post-id=1 -body="Nice post" -author="bob"
go run ./api-service/cmd/cli comments -post-id=1 -limit=20
go run ./api-service/cmd/cli search -query="Hello"
go run ./api-service/cmd/cli report-comments -limit=10
go run ./api-service/cmd/cli report-days
go run ./api-service/cmd/cli report-views -limit=10
```

## API Service endpoints

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

### Поиск и отчеты (через Data Service)

```bash
curl 'http://localhost:8080/api/v1/search?query=Hello'
curl 'http://localhost:8080/api/v1/reports/top-posts-by-comments?limit=10'
curl 'http://localhost:8080/api/v1/reports/posts-comments-by-day'
curl 'http://localhost:8080/api/v1/reports/top-posts-by-views?limit=10'
```

## Что реализовано по ТЗ

- 5 контейнеров: `api-service`, `data-service`, `postgres`, `zookeeper`, `kafka`
- `api-service` принимает HTTP и пишет в Kafka
- `data-service` читает Kafka и пишет в PostgreSQL
- минимум 2 связанные таблицы (есть `posts -> comments`, также `posts -> likes/views`)
- 3 отчета:
  - топ постов по количеству комментариев
  - количество постов и комментариев по дням
  - топ постов по просмотрам
