// handler.go — бизнес логика телеграм бота
package telegram

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageService - интерфейс для сервиса отправки приветсвенного сообщение пользователю
type MessageService interface {
	GetWelcomeMessage() string
}

// ActivityLogger - интерфейс то как мы хотим чтобы было логирование телеграмма
type ActivityLogger interface {
	LogTelegramUsersUse(userID int64, text string, duration time.Duration)
}

// Handler — это структура, которая знает, как отвечать на сообщения.
// В будущем мы добавим сюда services (Мозг), чтобы делать умные вещи.
type Handler struct {
	bot      *tgbotapi.BotAPI
	services MessageService
	logger   ActivityLogger
	commands map[string]func(*tgbotapi.Message)
}

// NewHandler
func NewHandler(bot *tgbotapi.BotAPI, services MessageService, logger ActivityLogger) *Handler {
	h := &Handler{
		bot:      bot,
		services: services,
		logger:   logger,
		commands: make(map[string]func(*tgbotapi.Message)),
	}
	h.initCommands()
	return h
}

func (h *Handler) initCommands() {
	h.commands["start"] = h.handleStart
}

// Handle - единая точка входа для обработки обновлений
func (h *Handler) Handle(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	// 1. засекаем время для логирования обработки сообщения
	start := time.Now()

	// Логика маршрутизации теперь здесь
	if update.Message.IsCommand() {
		command := update.Message.Command()
		if Handler, ok := h.commands[command]; ok {
			Handler(update.Message)
		} else {
			h.handleUnknown(update.Message)
		}
	} else {
		h.handleUnknown(update.Message)
	}

	// 2. отправляем счетчик прорабу
	h.logger.LogTelegramUsersUse(update.Message.Chat.ID, update.Message.Text, time.Since(start))
}

// handleStart - обрабатывает команду /start
func (h *Handler) handleStart(message *tgbotapi.Message) {
	// 1. Формируем текст ответа
	text := h.services.GetWelcomeMessage()
	// 2. Создаем сообщение для конкретного юзера
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	// 3. Отправляем сообщение
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("ошибка отправки сообщения: %v", err)
	}
}

// handleUnknown — реакция на неизвестную команду
func (h *Handler) handleUnknown(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды, введите /start")
	h.bot.Send(msg)
}
