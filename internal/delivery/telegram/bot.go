// bot.go — это файл, который отвечает за запуск бота, передача
// обновлений в обработчик
package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	handler *Handler
}

// запускаем бота
func NewBot(api *tgbotapi.BotAPI, handler *Handler) *Bot {
	return &Bot{
		api:     api,
		handler: handler,
	}
}

// Start запускает бесконечный цикл
func (b *Bot) Start() {
	log.Printf("Telegram bot run: %s", b.api.Self.UserName)
	// 1. создаем конфигурацию для получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// 2. подключаемся к Telegram и получаем канал, из которого будем читать все новые сообщения
	updates := b.api.GetUpdatesChan(u)

	// 3. цикл получения обновлений
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// 4. передаем сообщение в обработчик, он сам решит что с ним делать
		b.handler.Handle(update)
	}
}
