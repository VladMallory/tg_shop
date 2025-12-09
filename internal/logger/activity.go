package logger

import (
	"fmt"
	"io"
	"log"
	"time"
)

type Activity struct {
	UserID  any           // ID пользователя
	Event   string        // тип события
	Text    string        // текст
	Latency time.Duration // время выполнения
}

// ActivityLogger - писарь
// он хранит внутри себя инструменты для записи в файл
type ActivityLogger struct {
	loggers []*log.Logger // список всех мест куда надо писать
}

// NewActivityLogger принимает любое количество мест для записи
func NewActivityLogger(writers ...io.Writer) *ActivityLogger {
	var loggers []*log.Logger

	for _, w := range writers {
		loggers = append(loggers, log.New(w, "", log.LstdFlags))
	}

	return &ActivityLogger{
		loggers: loggers,
	}
}

func (l *ActivityLogger) Log(a Activity) {
	// 1. формируем сообщение один раз
	msg := fmt.Sprintf("[%v] [%s] [%v] %s", a.UserID, a.Event, a.Latency, a.Text)

	// 2. пишем в все места
	for _, logger := range l.loggers {
		logger.Println(msg)
	}
}
