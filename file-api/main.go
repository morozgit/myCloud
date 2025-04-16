package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

// Структура для хранения данных из сообщения
type Message struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

// Функция для скачивания файла
func downloadFile(filepath, destination string) error {
	// Проверяем, существует ли файл по указанному пути
	sourceFile, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer sourceFile.Close()

	// Создаем директорию для сохранения файла, если ее нет
	err = os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return fmt.Errorf("не удалось создать директорию: %w", err)
	}

	// Определяем путь для нового файла
	destinationFilePath := destination + "/" + filepath[strings.LastIndex(filepath, "/")+1:]

	// Создаем файл в новом месте
	destinationFile, err := os.Create(destinationFilePath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer destinationFile.Close()

	// Копируем содержимое из исходного файла в новый
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("не удалось скопировать данные: %w", err)
	}
	return nil
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
		"messages", // имя
		false,      // durable
		false,      // удалять, когда не используется
		false,      // эксклюзивный
		false,      // без ожидания
		nil,        // аргументы
	)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // очередь
		"",     // потребитель
		true,   // авто-подтверждение
		false,  // эксклюзивный
		false,  // не локальный
		false,  // без ожидания
		nil,    // аргументы
	)
	if err != nil {
		log.Fatalf("Не удалось зарегистрировать потребителя: %s", err)
	}

	fmt.Println("Ожидаем сообщения. Для выхода нажмите CTRL+C")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg Message
			// Декодируем сообщение JSON в структуру
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Printf("Ошибка при декодировании сообщения: %v", err)
				continue
			}

			// Извлекаем путь из сообщения
			filepath := msg.Path
			fmt.Printf("Получено сообщение с путем: %s\n", filepath)

			// Указываем, куда файл должен быть сохранен
			destinationPath := "/home/user/Downloads/MyCloudFiles"

			// Вызываем функцию для скачивания файла
			err = downloadFile(filepath, destinationPath)
			if err != nil {
				log.Printf("Ошибка при скачивании файла %s: %v", filepath, err)
			} else {
				fmt.Printf("Файл %s успешно скачан в %s\n", filepath, destinationPath)
			}
		}
	}()

	<-forever
}
