package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mycloud/file_api/filehandler"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/streadway/amqp"
)

type Message struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	queueName   = "file"
	serverAddr  = ":8080"
	filesRoute  = "/files/"
)

func main() {
	startFileServer()

	go handleMessages()

	fmt.Printf("Сервер запущен на http://localhost%s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

func startFileServer() {
	http.HandleFunc(filesRoute, func(w http.ResponseWriter, r *http.Request) {
		decodedPath, err := url.PathUnescape(r.URL.Path[len(filesRoute):])
		if err != nil {
			http.Error(w, "Неверный путь", http.StatusBadRequest)
			return
		}

		fullPath := filepath.Join("/home", decodedPath)
		fileName := filepath.Base(decodedPath)

		fileInfo, err := os.Stat(fullPath)
		if os.IsNotExist(err) || fileInfo.IsDir() {
			http.NotFound(w, r)
			return
		}

		file, err := os.Open(fullPath)
		if err != nil {
			http.Error(w, "Ошибка при открытии файла", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

		if _, err = io.Copy(w, file); err != nil {
			log.Printf("Ошибка при передаче файла: %v", err)
		}
	})
}

func handleMessages() {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось открыть канал: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Не удалось зарегистрировать потребителя: %v", err)
	}

	for d := range msgs {
		var msg Message
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			log.Printf("Ошибка при декодировании сообщения: %v", err)
			continue
		}

		fmt.Printf("Получено сообщение: путь=%s, имя=%s\n", msg.Path, msg.Name)

		link, err := filehandler.CreateDownloadLink(msg.Path)
		if err != nil {
			log.Printf("Ошибка при создании ссылки: %v", err)
			continue
		}

		fmt.Printf("Ссылка на скачивание: %s\n", link)

		if err := filehandler.DownloadFile(link); err != nil {
			log.Printf("Ошибка при загрузке файла %s: %v", msg.Name, err)
		} else {
			fmt.Printf("Файл %s успешно загружен\n", msg.Name)
		}
	}
}
