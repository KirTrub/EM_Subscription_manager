# ЭТАП 1: Сборка
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs
EXPOSE 3000
CMD ["./main"]