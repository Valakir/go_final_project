# Официальный образ Golang на базе Alpine как этап сборки
FROM golang:1.23-alpine AS builder

# Установка необходимых для сборки утилит
RUN apk add --no-cache gcc musl-dev

# Рабочая директория внутри контейнера
WORKDIR /app

# Копирование всех файлов в контейнер, чтобы иметь доступ к go.mod, go.sum и исходному коду
COPY . .

# Установка зависимостей
RUN go mod download

# Сборка бинарного файла
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o my_app cmd/main.go

# Используем минимальный образ на базе Alpine для финального контейнера
FROM alpine:latest

# Установка необходимых для выполнения библиотек
RUN apk add --no-cache ca-certificates

# Рабочая директория внутри контейнера
WORKDIR /app

# Копирование скомпилированного приложения из промежуточного контейнера
COPY --from=builder /app/my_app .

# Копирование директории web в финальный контейнер
COPY ./web ./web

# Определение переменных окружения, используемых приложением
ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=1234

# Порт, на котором будет работать приложение
EXPOSE ${TODO_PORT}

# Команда по умолчанию для запуска приложения
CMD ["/app/my_app"]