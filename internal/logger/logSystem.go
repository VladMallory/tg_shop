package logger

import (
	"salle_parfume/internal/logger/telegram"
	"time"
)

// LogSystem - главный прораб по логированию
type LogSystem struct {
}

// NewLogSystem - нанимаем прораба
func NewLogSystem() *LogSystem {
	return &LogSystem{}
}

// LogTelegram - логируем активность пользователя в телеграм
func (l *LogSystem) LogTelegram(userID int64, text string, duration time.Duration) {
	telegram.LogUserActivity(userID, text, duration)
}
