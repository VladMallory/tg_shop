package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"salle_parfume/internal/logger/telegram"
	"time"
)

// LogSystem - главный прораб по логированию
type LogSystem struct {
	telegramLogFile *os.File // хранение файла
}

// NewLogSystem - нанимаем прораба
func NewLogSystem() *LogSystem {
	// 1. путь где хранятся логи
	dir := "logs/telegram"
	filePath := filepath.Join(dir, "telegramUsersUse.txt")

	// 2. создаем папку если нету
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка создания папки для логов:", err)
		return nil
	}

	// 3. открываем файл при записи один раз при запуске программы
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка открытия файла для логов:", err)
		return &LogSystem{}
	}

	return &LogSystem{
		telegramLogFile: file,
	}
}

func (l *LogSystem) Close() {
	// существует ли файл для логов
	if l.telegramLogFile != nil {
		l.telegramLogFile.Close()
	}
}

// LogTelegram - логируем активность пользователя в телеграм
func (l *LogSystem) LogTelegram(userID int64, text string, duration time.Duration) {
	// 1. проверяем, что файл не nil
	if l.telegramLogFile == nil {
		fmt.Fprintln(os.Stderr, "ошибка: файл для логов не инициализирован")
		return
	}

	// 2. начинаем логирование активности пользователя в телеграм
	telegram.LogUserActivity(l.telegramLogFile, userID, text, duration)
}
