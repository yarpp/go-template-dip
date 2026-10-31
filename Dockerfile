# ---------- build ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Точка входа — ./api/main.go, бинарник кладём в /app/bin/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./api/main.go

# ---------- runtime ----------
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Копируем именно бинарник из /app/bin/api в /app/api
COPY --from=builder /app/bin/api /app/api

# Миграции нужны только для PostgreSQL-ветки; безопасно копировать
COPY --from=builder /app/migrations ./migrations

RUN mkdir -p /app/data
ENV DATA_DIR=/app/data
ENV LOGS_FILE=/app/log.txt

EXPOSE 8080
CMD ["/app/api"]