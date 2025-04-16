package main

import (
	"encoding/json"
	"fmt"
	"log"

	"mycloud/file_api/filehandler"

	"github.com/streadway/amqp"
)

type Message struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось открыть канал: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"messages",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Не удалось зарегистрировать потребителя: %s", err)
	}

	fmt.Println("Ожидаем сообщения. Для выхода нажмите CTRL+C")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg Message
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Printf("Ошибка при декодировании сообщения: %v", err)
				continue
			}

			filepath := msg.Path
			fmt.Printf("Получено сообщение с путем: %s\n", filepath)

			destinationPath := "/home/user/Downloads/MyCloudFiles"

			err = filehandler.DownloadFile(filepath, destinationPath)
			if err != nil {
				log.Printf("Ошибка при скачивании файла %s: %v", filepath, err)
			} else {
				fmt.Printf("Файл %s успешно скачан в %s\n", filepath, destinationPath)
			}
		}
	}()

	<-forever
}
