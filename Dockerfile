# Dockerfile
FROM golang:1.25.3-alpine

WORKDIR /app

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарник
RUN go build -o server ./cmd/server

# Порт приложения
EXPOSE 8080

# Запуск приложения
CMD ["./server"]
