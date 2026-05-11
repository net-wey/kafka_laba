FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/consumer ./cmd

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/consumer /app/consumer
EXPOSE 8081
ENTRYPOINT ["/app/consumer"]
