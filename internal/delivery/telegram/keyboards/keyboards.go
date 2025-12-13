package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Service отвечает за создание клавиатуры
type Service struct{}

func NewService() *Service {
	return &Service{}
}

// GetMainMenu возаращает главную клавиатуру после нажатия /start
func (s *Service) GetMainMenu() tgbotapi.InlineKeyboardMarkup {
	keyboards := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Каталог", "catalog"),
			tgbotapi.NewInlineKeyboardButtonData("О нас", "about"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Помощь", "help"),
		),
	)

	return keyboards
}
