# Официальный образ Golang на базе Alpine как этап сборки
FROM golang:1.22-alpine AS builder

# Установка необходимых для сборки утилит
RUN apk add --no-cache gcc musl-dev

#  Рабочая директорию внутри контейнера
WORKDIR /app

# Копирование go.mod и go.sum в контейнер и кеширование зависимости
COPY go.mod go.sum ./
RUN go mod download


# Копируем оставшуюся часть исходного кода в контейнер
COPY . .

# Сборка бинарного файла
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./my_app

# Используем минимальный образ на базе Alpine для финального контейнера
FROM alpine:latest

# Установка необходимых для выполнения библиотек
RUN apk add --no-cache ca-certificates

# Рабочая директорию внутри контейнера
WORKDIR /app

# Копируем скомпилированное приложение из промежуточного контейнера
COPY --from=builder /app/my_app /app/my_app
COPY ./web /app/web

# Определим переменные окружения, используемые приложением
ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=1234

#  Порт, на котором будет работать приложение
EXPOSE ${TODO_PORT}

# Команда по умолчанию для запуска приложения
CMD ["/app/my_app"]
