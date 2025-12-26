// app.go запускает всю программу. А так же занимается сборкой
package app

import (
	"fmt"
	"log"
	"salle_parfume/internal/config"
	"salle_parfume/internal/delivery/telegram"
	"salle_parfume/internal/delivery/telegram/keyboards"
	"salle_parfume/internal/logger"
	tgLogger "salle_parfume/internal/logger/telegram"
	"salle_parfume/internal/repository"
	"salle_parfume/internal/repository/sqlite"
	"salle_parfume/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// App - зависимости
type App struct {
	bot       *telegram.Bot     // telegram бот
	logWriter *logger.LogWriter // логгер
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
	logWriter, err := logger.NewLogWriter("telegram", "telegramUsersUse.txt")
	if err != nil {
		log.Fatal("ошибка в системе логирования", err)
	}
	activityLogger := tgLogger.NewLogger(logWriter)

	// 3. Инициализируем бота
	// Сначала создаем API
	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации API бота: %w", err)
	}

	// создаем сервис исходных сообщений
	messageService := service.NewMessageService()

	// создаем сервис клавиатур
	keyboardsService := keyboards.NewService()

	// 4. Инициализация db
	db, err := sqlite.NewSqliteDB(sqlite.Config{
		DriverName: "sqlite3",
		Path:       "./assets/storage.db",
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации DB: %w", err)
	}

	// инициализируем репозитории
	authRepo := sqlite.NewAuthSqlite(db)
	prodRepo := sqlite.NewProductSqlite(db)

	// собиаем все в один контейнер репозиториев
	repo := repository.NewRepository(authRepo, prodRepo)

	// Создаем Handler (он принимает API и Сервис сообщений)
	handler := telegram.NewHandler(botAPI, messageService, activityLogger, keyboardsService, repo, cfg.AdminID)

	// Создаем самого бота (принимает API, Handler)
	bot := telegram.NewBot(botAPI, handler)

	// возвращаем готового, сборанного приложения
	return &App{
		bot:       bot,
		logWriter: logWriter,
	}, nil
}

// Run сценарий работы. Этап запуска приложения
// бизнес логика, без мусора if else и прочее
func (a *App) Run() {
	// закрываем файл логов
	defer a.logWriter.Close()

	// запускаем бота
	a.bot.Start()
}
