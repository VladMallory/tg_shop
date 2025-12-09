// задача: прочиать .env и вернуть значения
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
}

func LoadConfig() (*Config, error) {
	// 1. Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env пуст или не найден")
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
