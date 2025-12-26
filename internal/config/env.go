// env.go читает .env и возвращает данные из переменных
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string // для работы телеграм бота
	AdminID       int64  // админу будет достопно добавление нового каталога
}

func LoadConfig() (*Config, error) {
	// Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env пуст или файл не создан в дериктории")
	}

	// 1. Достаем токен TELEGRAM_TOKEN
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("укажите в env TELEGRAM_TOKEN")
	}

	// 2. Достаем TELEGRAM_ID
	adminIDStr := os.Getenv("TELEGRAM_ID")
	if adminIDStr == "" {
		return nil, fmt.Errorf("укажите в .env TELEGRAM_ID")
	}

	// преобразуем строку в число
	adminIDInt, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("не получилось преобразовать ADMIN_ID")
	}

	return &Config{
		TelegramToken: token,
		AdminID:       adminIDInt,
	}, nil
}
