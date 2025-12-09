package app

import (
	"fmt"
	"os"
	"salle_parfume/internal/config"
	"salle_parfume/internal/delivery/telegram"
	"salle_parfume/internal/logger"
	"salle_parfume/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// App - зависимости
type App struct {
	bot     *telegram.Bot          // telegram бот
	logger  *logger.ActivityLogger // наш логер
	logFile *os.File               // файл логов (теперь мы храним его здесь)
}

// New эта сборки. Конструктор
// тут все проверки ошибок
func New() (*App, error) {
	// 1. Грузим env
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфига: %w", err)
	}

	// 2. Настраиваем логирование
	activityLogger, logFile, err := logger.SystemLogs()
	if err != nil {
		return nil, fmt.Errorf("ошибка настройки логов: %w", err)
	}

	// 3. Инициализируем бота
	// Сначала создаем API
	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации API бота: %w", err)
	}

	// создаем сервис исходных сообщений
	messageService := service.NewMessageService()

	// Создаем Handler (он принимает API и Сервис сообщений)
	handler := telegram.NewHandler(botAPI, messageService)

	// Создаем самого бота (принимает API, Handler и Logger)
	bot := telegram.NewBot(botAPI, handler, activityLogger)

	// возвращаем готового, сборанного приложения
	return &App{
		bot:     bot,
		logger:  activityLogger,
		logFile: logFile,
	}, nil
}

// Run сценарий работы. Этап запуска приложения
// бизнес логика, без мусора if else и прочее
func (a *App) Run() {
	if a.logFile != nil {
		defer a.logFile.Close()
	}
	// запускаем бота
	a.bot.Start()
}
