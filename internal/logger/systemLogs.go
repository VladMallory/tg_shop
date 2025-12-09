package logger

import "os"

// SystemLogs - оркестратор логов
// Он запускает все системы логирования
func SystemLogs() (*ActivityLogger, *os.File, error) {
	// 1. Запускаем логирование телеграма
	// Мы обязаны вернуть логер, чтобы передать его боту
	tgLogger, tgFile, err := SetupTelegram()
	if err != nil {
		return nil, nil, err
	}

	// В будущем тут будет:
	// siteLogger := SetupSite()
	// appLogger := SetupApp()

	return tgLogger, tgFile, nil
}
