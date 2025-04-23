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
	"strings"

	"github.com/streadway/amqp"
)

type Message struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Content string `json:"file_data"`
}

const (
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
	go handleDownloadQueue()
	go handleUploadQueue()

}

func getRabbitMQChannel() (*amqp.Channel, *amqp.Connection, error) {
	conn, err := amqp.Dial(config.GetRabbitMQURL())
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось подключиться к RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("не удалось открыть канал RabbitMQ: %w", err)
	}

	return ch, conn, nil
}

func handleDownloadQueue() {
	ch, conn, err := getRabbitMQChannel()
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare("file", false, false, false, false, nil)
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

		if msg.URL != "" {
			log.Println("Файл уже был обработан или ссылка отправлена")
			continue
		}

		msg.URL = link
		sendMessage(msg)
	}
}

func sendMessage(msg Message) {
	ch, conn, err := getRabbitMQChannel()
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	_, err = ch.QueueDeclare("get_link", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %v", err)
	}

	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Ошибка при маршаллинге сообщения: %v", err)
		return
	}

	err = ch.Publish("", "get_link", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        encodedMsg,
	})
	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
		return
	}

	fmt.Printf("Ссылка на скачивание отправлена в очередь: %s\n", msg.URL)
}

func resolveUploadPath(originalPath string) (string, error) {
	rootDir := os.Getenv("UPLOAD_ROOT_DIR")
	if rootDir == "" {
		return "", fmt.Errorf("UPLOAD_ROOT_DIR не установлен в .env")
	}

	// Приводим путь к относительному, если он начинается с /
	relativePath := strings.TrimPrefix(originalPath, "/")

	// Склеиваем с корневым путем
	finalPath := filepath.Join(rootDir, relativePath)
	return finalPath, nil
}

func handleUploadQueue() {
	ch, conn, err := getRabbitMQChannel()
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare("upload", false, false, false, false, nil)
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

		resolvedPath, err := resolveUploadPath(msg.Path)
		if err != nil {
			log.Printf("Ошибка при разрешении пути: %v", err)
			continue
		}
		msg.Path = resolvedPath

		if err := os.MkdirAll(msg.Path, 0755); err != nil {
			log.Printf("Ошибка при создании директории: %v", err)
			continue
		}

		filePath := filepath.Join(msg.Path, msg.Name)
		err = os.WriteFile(filePath, []byte(msg.Content), 0644)
		if err != nil {
			log.Printf("Ошибка при сохранении файла: %v", err)
			continue
		}

		fmt.Printf("Файл успешно сохранен: %s\n", filePath)

	}
}
