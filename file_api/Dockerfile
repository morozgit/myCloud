FROM golang:1.23-alpine

WORKDIR /file_api

COPY . .

RUN go mod tidy && go build -o file_api .

# Добавляем netcat
RUN apk add --no-cache netcat-openbsd

CMD sh -c "while ! nc -z rabbitmq 5672; do echo '⏳ Жду RabbitMQ...'; sleep 1; done && ./file_api"
