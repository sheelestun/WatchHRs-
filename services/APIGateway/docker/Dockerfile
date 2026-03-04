# Build stage
FROM golang:1.25-alpine AS builder

# Установка зависимостей для сборки
RUN apk add --no-cache git ca-certificates

# Создание рабочей директории
WORKDIR /build

# Копирование go.mod и go.sum для кеширования зависимостей
COPY APIGateway/go.mod APIGateway/go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN cd APIGateway && go build -o api-gateway ./cmd/main

# Production stage
FROM alpine:latest

# Установка необходимых утилит
RUN apk --no-cache add ca-certificates tzdata

# Создание пользователя для безопасности
RUN adduser -D -s /bin/sh appuser

# Создание рабочей директории
WORKDIR /app

# Создание директории для логов
RUN mkdir -p /app/logs && chown appuser:appuser /app/logs

# Копирование бинарного файла из builder stage
COPY --from=builder /build/APIGateway/api-gateway .

# Изменение владельца файлов
RUN chown appuser:appuser api-gateway

# Переключение на пользователя appuser
USER appuser

# Экспорт порта
EXPOSE 8080

# Запуск приложения
CMD ["./api-gateway"]