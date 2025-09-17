FROM golang:1.24.3 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Статическая сборка бинарника
#RUN go build -o video-manage ./cmd/video-manage
RUN CGO_ENABLED=0 GOOS=linux go build -o sub-manage ./cmd/sub-manage
FROM debian:bullseye-slim

WORKDIR /app

# Установка системных библиотек (для SSL, DNS, логов и т.п.)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Копируем из builder-а бинарь, конфиги и миграции
COPY --from=builder /app/sub-manage .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/internal/migrations /app/migrations

# Открываем порт для gRPC (опционально, для понимания)
#EXPOSE 50051

CMD ["./sub-manage"]
