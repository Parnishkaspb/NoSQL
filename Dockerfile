FROM golang:1.23.2

WORKDIR /app

COPY . .

RUN #go build -o server ./cmd/main.go

COPY wait-and-start.sh ./wait-and-start.sh
RUN chmod +x ./wait-and-start.sh

CMD ["./wait-and-start.sh"]