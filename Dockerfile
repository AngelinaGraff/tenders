# Используем базовый образ с Go версии 1.20
FROM golang:1.20-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum файлы
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем весь код в контейнер
COPY . .

# Собираем приложение
RUN go build -o server cmd/server/main.go

# Указываем переменные окружения
ENV PORT=8080
ENV GIN_MODE=release

# Открываем порт для приложения
EXPOSE 8080

# Команда для запуска приложения
CMD ["./server"]