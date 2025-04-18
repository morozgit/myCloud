package rabbithandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mycloud/file_api/config"
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
	URL  string `json:"url"`
}

const (
	queueName  = "file"
	filesRoute = "/files/"
)

func StartFileServer() {
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

func HandleMessages() {
	var rabbitMQURL = config.GetRabbitMQURL()
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

	// Обрабатываем сообщения
	for d := range msgs {
		var msg Message
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			log.Printf("Ошибка при декодировании сообщения: %v", err)
			continue
		}

		fmt.Printf("Получено сообщение: путь=%s, имя=%s\n", msg.Path, msg.Name)

		// Создание ссылки для скачивания
		link, err := filehandler.CreateDownloadLink(msg.Path)
		if err != nil {
			log.Printf("Ошибка при создании ссылки: %v", err)
			continue
		}

		// Если ссылка уже существует, пропускаем обработку
		if msg.URL != "" {
			log.Println("Файл уже был обработан или ссылка отправлена")
			continue
		}

		// Обновляем URL
		msg.URL = link

		// Отправляем сообщение обратно в очередь для дальнейшей обработки
		// В данном случае вызовем функцию для отправки сообщения
		sendMessage(msg)
	}
}

func sendMessage(msg Message) {
	var rabbitMQURL = config.GetRabbitMQURL()
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

	// Просто вызываем QueueDeclare без присваивания переменной, если не нужно
	_, err = ch.QueueDeclare("get_link", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %v", err)
	}

	// Кодируем сообщение в формат JSON
	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Ошибка при маршаллинге сообщения: %v", err)
		return
	}

	// Отправляем сообщение в очередь
	err = ch.Publish("", "get_link", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        encodedMsg,
	})
	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
		return
	}

	// Логируем успешную отправку сообщения
	fmt.Printf("Ссылка на скачивание отправлена в очередь: %s\n", msg.URL)
}
