// handler.go — бизнес логика телеграм бота
package telegram

import (
	"log"
	"salle_parfume/internal/logger"
	"salle_parfume/internal/service"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler — это структура, которая знает, как отвечать на сообщения.
// В будущем мы добавим сюда services (Мозг), чтобы делать умные вещи.
type Handler struct {
	bot       *tgbotapi.BotAPI
	services  *service.MessageService // сервис для обработки исходящих сообщений для пользователя
	logSystem *logger.LogSystem
}

// NewHandler — конструктор
func NewHandler(bot *tgbotapi.BotAPI, services *service.MessageService, logSystem *logger.LogSystem) *Handler {
	return &Handler{
		bot:       bot,
		services:  services,
		logSystem: logSystem,
	}
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
		switch update.Message.Command() {
		case "start":
			h.handleStart(update.Message)
		default:
			h.handleUnknown(update.Message)
		}
	} else {
		h.handleMessage(update.Message)
	}

	// 2. отправляем счетчик прорабу
	h.logSystem.LogTelegram(update.Message.Chat.ID, update.Message.Text, time.Since(start))
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
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды")
	h.bot.Send(msg)
}

// handleMessage — реакция на просто текст (не команду)
func (h *Handler) handleMessage(message *tgbotapi.Message) {
	// 1. Создаем сообщение
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я пока понимаю только /start")
	// 2. Отправляем сообщение
	h.bot.Send(msg)
}
