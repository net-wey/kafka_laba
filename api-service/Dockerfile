FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api-service ./cmd

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/api-service /app/api-service
EXPOSE 8080
ENTRYPOINT ["/app/api-service"]
