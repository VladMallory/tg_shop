package telegram

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogUserActivity - логируем активность пользователя в телеграм
func LogUserActivity(userID int64, text string, duration time.Duration) {
	// 1. формируем строку
	logLine := fmt.Sprintf("[%d] [%dms] %s\n", userID, duration.Milliseconds(), text)

	// 2. определяем путь лога
	dir := "logs/telegram"
	filePath := filepath.Join(dir, "telegramUsersUse.txt")

	// 3. гарантируем что папка есть
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка создания папки логов для telegram:", err)
		return
	}

	// 4. открываем файл
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка открытия файла логов для telegram:", err)
		return
	}
	defer file.Close()

	// 5. пишем в файл
	if _, err := file.WriteString(logLine); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка записи в файл логов для telegram:", err)
		return
	}

}
