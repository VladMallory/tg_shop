// env.go читает .env и возвращает данные из переменных
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
}

func LoadConfig() (*Config, error) {
	// 1. Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env пуст или файл не создан в дериктории")
	}

	// 2. Достаем токен
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("укажите в env TELEGRAM_TOKEN")
	}

	// 3. Возвращаем конфиг в структуру
	return &Config{
		TelegramToken: token,
	}, nil
}
