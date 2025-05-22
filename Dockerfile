FROM golang:1.23.2

WORKDIR /app

# Копируем весь проект
COPY . .

# Собираем Go-приложение
RUN #go build -o server ./cmd/main.go

# Копируем скрипт ожидания MongoDB
COPY wait-and-start.sh ./wait-and-start.sh
RUN chmod +x ./wait-and-start.sh

# Используем скрипт как команду запуска
CMD ["./wait-and-start.sh"]