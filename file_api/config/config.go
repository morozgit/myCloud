package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_ = godotenv.Load()
}

func GetRabbitMQURL() string {
	return os.Getenv("RABBITMQ_URL")
}

func GetBaseURL() string {
	return os.Getenv("BASE_URL")
}
