// logSystem.go - общая система логирования в файлы
package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

// LogWriter - главный прораб по логированию
type LogWriter struct {
	file *os.File // передаем файл логов
}

// NewLogWriter - нанимаем прораба
func NewLogWriter(subDir, fileName string) (*LogWriter, error) {
	// 1. определяем путь где хранятся логи
	dir := filepath.Join("logs", subDir)

	// 2. то как будет в итоге называться файл
	filePath := filepath.Join(dir, fileName)

	// 3. создаем папку если нету
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("ошибка создания папки для логов: %v", err)
	}

	// 4. открываем файл при записи один раз при запуске программы
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка открытия файла для логов:", err)
		return &LogWriter{
			file: file,
		}, nil
	}

	return &LogWriter{
		file: file,
	}, nil
}

// позволяет закрыть файл в app.go
func (l *LogWriter) Close() {
	// существует ли файл для логов
	if l.file != nil {
		l.file.Close()
	}
}

// Write - пишет строку в файл
func (l *LogWriter) Write(text string) {
	// 1. проверяем, что файл не nil
	if l.file == nil {
		fmt.Fprintln(os.Stderr, "ошибка: файл для логов не инициализирован")
		return
	}

	if _, err := l.file.WriteString(text); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка записи в файл логов", err)
	}
}
