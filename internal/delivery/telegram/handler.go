// handler.go — бизнес логика телеграм бота
package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"salle_parfume/internal/domain"
	"salle_parfume/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageService - интерфейс для сервиса отправки приветсвенного сообщение пользователю
type MessageService interface {
	GetWelcomeMessage() string
	GetAboutMessage() string
}

// ActivityLogger - интерфейс то как мы хотим чтобы было логирование телеграмма
type ActivityLogger interface {
	LogTelegramUsersUse(userID int64, text string, duration time.Duration)
}

// KeyboardProvider - интерфейс для предоставления клавиатур
type KeyboardProvider interface {
	GetMainMenu() tgbotapi.InlineKeyboardMarkup
	GetProductTypeKeyboard() tgbotapi.InlineKeyboardMarkup // Добавили новый метод
	GetBuyKeyboard(productID int64) tgbotapi.InlineKeyboardMarkup
}

// Состояния FSM (Finite State Machine)
// Это этапы нашего диалога
type State int

const (
	StateNone                  State = iota
	StateWaitingForType              // Ждем выбор типа
	StateWaitingForPhoto             // Ждем фото
	StateWaitingForName              // Ждем название
	StateWaitingForDescription       // Ждем описание
	StateWaitingForPrice             // Ждем цену
)

// DraftProduct - временная структура (черновик), пока мы собираем данные
type DraftProduct struct {
	Type        domain.ProductType
	ImageID     string
	Name        string
	Description string
	Price       float64
}

// Handler — это структура, которая знает, как отвечать на сообщения.
type Handler struct {
	bot       *tgbotapi.BotAPI
	services  MessageService
	logger    ActivityLogger
	keyboards KeyboardProvider
	repo      repository.ProductRepository // Используем интерфейс, а не конкретную структуру
	adminID   int64                        // ID админа
	commands  map[string]func(*tgbotapi.Message)

	// Состояние пользователя (где он сейчас в диалоге)
	userStates map[int64]State
	// Черновики товаров для каждого пользователя
	drafts map[int64]*DraftProduct
}

// NewHandler создает новый обработчик
// Теперь принимает репозиторий (интерфейс) и ID админа
func NewHandler(bot *tgbotapi.BotAPI, services MessageService, logger ActivityLogger, keyboards KeyboardProvider, repo repository.ProductRepository, adminID int64) *Handler {
	h := &Handler{
		bot:        bot,
		services:   services,
		logger:     logger,
		keyboards:  keyboards,
		repo:       repo,
		adminID:    adminID,
		commands:   make(map[string]func(*tgbotapi.Message)),
		userStates: make(map[int64]State),
		drafts:     make(map[int64]*DraftProduct),
	}
	h.initCommands()
	return h
}

// initCommands инициализирует карту доступных команд бота
func (h *Handler) initCommands() {
	h.commands["start"] = h.handleStart
	h.commands["new"] = h.handleNewProduct
}

// Handle - единая точка входа для обработки обновлений
func (h *Handler) Handle(update tgbotapi.Update) {
	start := time.Now()

	// является ли это кнопкой. Если нет, пропускаем
	if update.CallbackQuery != nil {
		h.handleCallback(update.CallbackQuery)
		return
	}

	if update.Message == nil {
		return
	}

	// Проверяем, находится ли пользователь в процессе диалога
	if state, ok := h.userStates[update.Message.Chat.ID]; ok && state != StateNone {
		h.handleState(update.Message, state)
		return
	}

	// Обычная обработка команд
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

	h.logger.LogTelegramUsersUse(update.Message.Chat.ID, update.Message.Text, time.Since(start))
}

// handleNewProduct - начало процесса добавления товара
func (h *Handler) handleNewProduct(message *tgbotapi.Message) {
	// Проверка прав доступа (Security)
	if message.From.ID != h.adminID {
		h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "У вас нет прав для этой команды."))
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите тип духов:")
	msg.ReplyMarkup = h.keyboards.GetProductTypeKeyboard()
	h.bot.Send(msg)

	// Переводим пользователя в состояние "Ждем выбор типа"
	h.userStates[message.Chat.ID] = StateWaitingForType
	h.drafts[message.Chat.ID] = &DraftProduct{}
}

// handleCallback - обработка нажатий на кнопки
func (h *Handler) handleCallback(callback *tgbotapi.CallbackQuery) {
	// записываем телеграм id клиента
	chatID := callback.Message.Chat.ID
	// кнопка которая была нажата
	data := callback.Data

	log.Printf("Callback: chatID=%d, data=%s", chatID, data)

	// кнопка "Каталог"
	if data == "catalog" {
		h.handleCatalog(chatID)
		h.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		return
	}

	if data == "about" {
		h.handleAbout(chatID)
		h.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		return
	}

	// Обработка кнопки "Купить"
	if strings.HasPrefix(data, "buy_") {
		h.handleBuy(chatID, data)
		h.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		return
	}

	// Проверяем, если это выбор типа, но стейт потерян (например, после рестарта бота)
	if strings.HasPrefix(data, "type_") {
		state, ok := h.userStates[chatID]
		if !ok || state != StateWaitingForType {
			log.Printf("State mismatch or lost context for user %d. State: %v", chatID, state)
			h.bot.Send(tgbotapi.NewMessage(chatID, "Произошла ошибка (возможно, бот был перезапущен). Пожалуйста, введите /new заново."))
			h.bot.Request(tgbotapi.NewCallback(callback.ID, "")) // Убираем часики
			return
		}
	}

	// Если ждем тип духов
	if state, ok := h.userStates[chatID]; ok && state == StateWaitingForType {
		draft := h.drafts[chatID]
		if draft == nil {
			// Если вдруг драфта нет (хотя стейт есть - странно, но подстрахуемся)
			h.bot.Send(tgbotapi.NewMessage(chatID, "Внутренняя ошибка. Пожалуйста, начните заново /new"))
			h.userStates[chatID] = StateNone
			return
		}

		switch data {
		case "type_female":
			draft.Type = domain.TypeFemale
		case "type_male":
			draft.Type = domain.TypeMale
		case "type_unisex":
			draft.Type = domain.TypeUnisex
		default:
			return
		}

		log.Printf("User %d selected type: %s", chatID, draft.Type)

		// Переходим к следующему шагу
		h.userStates[chatID] = StateWaitingForPhoto
		msg := tgbotapi.NewMessage(chatID, "Отправьте фотографию:")
		h.bot.Send(msg)
		h.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
	}
}

// handleState - пошаговая обработка ввода данных
func (h *Handler) handleState(message *tgbotapi.Message, state State) {
	chatID := message.Chat.ID
	draft := h.drafts[chatID]

	switch state {
	case StateWaitingForPhoto:
		if message.Photo == nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Пожалуйста, отправьте фото."))
			return
		}
		// Берем самое качественное фото (последнее в массиве)
		photo := message.Photo[len(message.Photo)-1]
		draft.ImageID = photo.FileID

		h.userStates[chatID] = StateWaitingForName
		h.bot.Send(tgbotapi.NewMessage(chatID, "Введите название:"))

	case StateWaitingForName:
		draft.Name = message.Text
		h.userStates[chatID] = StateWaitingForDescription
		h.bot.Send(tgbotapi.NewMessage(chatID, "Введите описание:"))

	case StateWaitingForDescription:
		draft.Description = message.Text
		h.userStates[chatID] = StateWaitingForPrice
		h.bot.Send(tgbotapi.NewMessage(chatID, "Введите цену товара:"))

	case StateWaitingForPrice:
		price, err := strconv.ParseFloat(message.Text, 64)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Пожалуйста, введите корректное число."))
			return
		}
		draft.Price = price

		// Сохраняем готовый товар в базу
		product := &domain.Product{
			Type:        draft.Type,
			Name:        draft.Name,
			Description: draft.Description,
			Price:       draft.Price,
			ImageID:     draft.ImageID,
		}

		if err := h.repo.CreateProduct(product); err != nil {
			log.Printf("Error creating product: %v", err)
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при сохранении товара."))
		} else {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Готово, духи добавлены в каталог!"))
		}

		// Сбрасываем состояние
		h.userStates[chatID] = StateNone
		delete(h.drafts, chatID)
	}
}

func (h *Handler) handleCatalog(chatID int64) {
	products, err := h.repo.GetAllProducts()
	if err != nil {
		log.Printf("Error getting products: %v", err)
		h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при получении каталога."))
		return
	}

	if len(products) == 0 {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Каталог пока пуст."))
		return
	}

	for _, p := range products {
		text := fmt.Sprintf("<b>%s</b>\n\n%s\n\nЦена: %.2f руб.", p.Name, p.Description, p.Price)
		msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(p.ImageID))
		msg.Caption = text
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = h.keyboards.GetBuyKeyboard(p.ID)
		h.bot.Send(msg)
	}
}

// handleBuy - обработка нажатия кнопки "Купить"
func (h *Handler) handleBuy(chatID int64, data string) {
	productID := strings.TrimPrefix(data, "buy_")
	h.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Вы выбрали товар с ID: %s. \nФункционал оформления заказа находится в разработке.", productID)))
}

func (h *Handler) handleAbout(chatID int64) {
	text := h.services.GetAboutMessage()
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

}

// handleStart - обрабатывает команду /start
func (h *Handler) handleStart(message *tgbotapi.Message) {
	// 1. Формируем текст ответа
	text := h.services.GetWelcomeMessage()
	// 2. Создаем сообщение для конкретного юзера
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	// 3. добавляем клавитуру к /start
	msg.ReplyMarkup = h.keyboards.GetMainMenu()

	// 4. Отправляем сообщение
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("ошибка отправки сообщения: %v", err)
	}
}

// handleUnknown — реакция на неизвестную команду
func (h *Handler) handleUnknown(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды, введите /start")
	h.bot.Send(msg)
}
