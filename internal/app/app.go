package app

import (
	"fmt"
	"salle_parfume/internal/config"
	"salle_parfume/internal/delivery/telegram"
	"salle_parfume/internal/logger"
	"salle_parfume/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// App - зависимости
type App struct {
	bot       *telegram.Bot     // telegram бот
	logSystem *logger.LogSystem // логгер
}

// New эта сборки. Конструктор
// тут все проверки ошибок
func New() (*App, error) {
	// 1. Грузим env
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфига: %w", err)
	}

	// 2. логирование
	logSystem := logger.NewLogSystem()

	// 3. Инициализируем бота
	// Сначала создаем API
	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации API бота: %w", err)
	}

	// создаем сервис исходных сообщений
	messageService := service.NewMessageService()

	// Создаем Handler (он принимает API и Сервис сообщений)
	handler := telegram.NewHandler(botAPI, messageService, logSystem)

	// Создаем самого бота (принимает API, Handler)
	bot := telegram.NewBot(botAPI, handler)

	// возвращаем готового, сборанного приложения
	return &App{
		bot:       bot,
		logSystem: logSystem,
	}, nil
}

// Run сценарий работы. Этап запуска приложения
// бизнес логика, без мусора if else и прочее
func (a *App) Run() {
	// закрываем файл логов
	defer a.logSystem.Close()

	// запускаем бота
	a.bot.Start()
}
