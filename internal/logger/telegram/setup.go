package telegram

import (
	"fmt"
	"os"
	"path/filepath"
	"salle_parfume/internal/logger"
)

// Setup - готовит логер специально для ТЕЛЕГРАМА
// Это специализация. Она знает путь к файлу.
func Setup() (*logger.ActivityLogger, *os.File, error) {
	const logPath = "logs/usersUse.txt"

	// 1. Создаем директорию
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, nil, fmt.Errorf("ошибка создания директории %s: %w", dir, err)
	}

	// 2. Открываем файл
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка открытия файла %s: %w", logPath, err)
	}

	// 3. Используем ОБЩИЙ логер из пакета-родителя
	// Передаем ему файл + консоль
	l := logger.NewActivityLogger(file, os.Stdout)

	return l, file, nil
}
