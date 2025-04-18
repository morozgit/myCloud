package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_ = godotenv.Load()
}

func GetRabbitMQURL() string {
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
}

func GetBaseURL() string {
	return os.Getenv("BASE_URL")
}
