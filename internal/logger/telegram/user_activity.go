package telegram

import (
	"fmt"
	"os"
	"time"
)

// LogUserActivity - логируем активность пользователя в телеграм (пишем в уже открытый файл)
func LogUserActivity(file *os.File, userID int64, text string, duration time.Duration) {
	// 1. формируем строку
	logLine := fmt.Sprintf("[%d] [%dms] %s\n", userID, duration.Milliseconds(), text)

	// 2. пишем в файл
	if _, err := file.WriteString(logLine); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка записи в файл логов для telegram:", err)
	}
}
