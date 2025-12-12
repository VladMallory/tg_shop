package telegram

import (
	"fmt"
	"salle_parfume/internal/logger"
	"time"
)

type Logger struct {
	writer *logger.LogWriter
}

// создаем логгер телеграмма
func NewLogger(writer *logger.LogWriter) *Logger {
	return &Logger{
		writer: writer,
	}
}

// LogUserActivity - логируем активность пользователя в телеграм (пишем в уже открытый файл)
func (l *Logger) LogTelegramUsersUse(userID int64, text string, duration time.Duration) {
	// 1. формируем строку
	logLine := fmt.Sprintf("[%d] [%dms] %s\n", userID, duration.Milliseconds(), text)

	// 2. пишем через помощника
	l.writer.Write(logLine)
}
