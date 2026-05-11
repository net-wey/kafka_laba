Как показать, что запись не напрямую в БД, а через Kafka:


```
docker compose exec postgres psql -U postgres -d kafka_laba -tAc "SELECT COUNT(*) FROM posts;"
```

```
docker compose stop consumer
```
```
curl -X POST http://localhost:8080/api/v1/posts \
  -H 'Content-Type: application/json' \
  -d '{"title":"Post 1","body":"Hello","author":"alice"}'
```
```
make cli ARGS='post -title="Kafka Proof" -body="queued" -author="proof"'
```

```
docker compose exec postgres psql -U postgres -d kafka_laba -tAc "SELECT COUNT(*) FROM posts;"
```

```
docker compose start consumer
sleep 5
```

```
docker compose exec postgres psql -U postgres -d kafka_laba -tAc "SELECT COUNT(*) FROM posts;"
```
