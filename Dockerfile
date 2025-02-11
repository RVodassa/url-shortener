# Этап сборки
FROM golang:1.23.3-alpine AS builder

WORKDIR /app/url-shortener

# Установка необходимых пакетов для сборки
RUN apk add --no-cache gcc musl-dev

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники приложения в рабочую директорию
COPY . .

# Компилируем приложение
RUN go build -o url-shortener ./main.go

# Этап выполнения
FROM alpine:latest

WORKDIR /app/url-shortener

# Копируем скомпилированное приложение из образа builder
COPY --from=builder /app/url-shortener/url-shortener .

# Копируем папку configs
COPY ./configs /app/url-shortener/configs
COPY ./migrations /app/url-shortener/migrations

COPY .env ./

# Открываем порт 8083
EXPOSE 8083

# Запускаем приложение
ENTRYPOINT ["sh", "-c", "./url-shortener -cfg_path $CFG_PATH -storage $STORAGE_TYPE"]





