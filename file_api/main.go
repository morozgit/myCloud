package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mycloud/file_api/filehandler"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

type Message struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func main() {
	go func() {
		fs := http.FileServer(http.Dir("/home"))
		http.Handle("/files/", http.StripPrefix("/files/", fs))
		log.Println("Файловый сервер запущен на http://0.0.0.0:8080")
		log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
	}()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла")
	}
	baseDir := os.Getenv("BASE_DIR")
	if baseDir == "" {
		log.Fatal("BASE_DIR не задан в .env")
	}
	downloadDir := baseDir + "Downloads/MyCloudFiles/"

	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось открыть канал: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("file", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %s", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
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
			filename := msg.Name
			fmt.Printf("Получено сообщение с путем: %s\n", filepath)
			fmt.Printf("Получено сообщение с именем: %s\n", filename)

			link, err := filehandler.CreateDownloadLink(filepath)
			if err != nil {
				log.Printf("Ошибка при создании ссылки %s: %v", filepath, err)
				continue
			}
			fmt.Printf("Ссылка на скачивание: %s\n", link)

			err = filehandler.DownloadFile(link, filename, downloadDir)
			if err != nil {
				log.Printf("Ошибка при загрузке файла %s: %v", filename, err)
			} else {
				fmt.Printf("Файл %s успешно загружен в папку %s\n", filename, downloadDir)
			}
		}
	}()

	<-forever
}
