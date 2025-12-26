// keyboards.go - Пакет keyboards отвечает за создание и управление клавиатурами (Inline и Reply) для Telegram бота.
// Он содержит методы для генерации различных меню, необходимых для навигации пользователя.
package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Константы для Callback Data кнопок.
// Использование констант помогает избежать ошибок "magic strings" и упрощает поддержку кода.
const (
	ButtonCatalog = "catalog"
	ButtonAbout   = "about"
	ButtonHelp    = "help"

	TypeFemale = "type_female"
	TypeMale   = "type_male"
	TypeUnisex = "type_unisex"

	PrefixBuy = "buy_%d"
)

// Service реализует логику создания клавиатур.
// Он является поставщиком (Provider) разметки для сообщений бота.
type Service struct{}

// NewService создает новый экземпляр сервиса клавиатур.
// Возвращает указатель на структуру Service.
func NewService() *Service {
	// В данном случае инициализация простая, но в будущем здесь может быть
	// внедрение зависимостей (например, локализация).
	return &Service{}
}

// GetMainMenu формирует и возвращает главную Inline-клавиатуру.
// Обычно отображается после команды /start.
func (s *Service) GetMainMenu() tgbotapi.InlineKeyboardMarkup {
	// Создаем клавиатуру с помощью tgbotapi
	keyboards := tgbotapi.NewInlineKeyboardMarkup(
		// Первый ряд кнопок
		tgbotapi.NewInlineKeyboardRow(
			// Кнопка "Каталог" отправляет callback_data "catalog"
			tgbotapi.NewInlineKeyboardButtonData("Каталог", ButtonCatalog),
			// Кнопка "О нас" отправляет callback_data "about"
			tgbotapi.NewInlineKeyboardButtonData("О нас", ButtonAbout),
		),
		// Второй ряд кнопок
		tgbotapi.NewInlineKeyboardRow(
			// Кнопка "Помощь" отправляет callback_data "help"
			tgbotapi.NewInlineKeyboardButtonData("Помощь", ButtonHelp),
		),
	)

	return keyboards
}

// GetProductTypeKeyboard создает клавиатуру для выбора категории товара.
// Используется администратором при добавлении новой позиции.
func (s *Service) GetProductTypeKeyboard() tgbotapi.InlineKeyboardMarkup {
	// Формируем разметку клавиатуры
	return tgbotapi.NewInlineKeyboardMarkup(
		// Ряд 1: Основные гендерные типы
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Женские", TypeFemale),
			tgbotapi.NewInlineKeyboardButtonData("Мужские", TypeMale),
		),
		// Ряд 2: Универсальный тип
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Унисекс", TypeUnisex),
		),
	)
}

// GetBuyKeyboard генерирует клавиатуру действия для конкретного товара.
// Принимает productID для формирования уникального callback_data.
func (s *Service) GetBuyKeyboard(productID int64) tgbotapi.InlineKeyboardMarkup {
	// Формируем строку callback_data с ID товара (например, "buy_123")
	callbackData := fmt.Sprintf(PrefixBuy, productID)

	// Создаем клавиатуру с одной кнопкой "Купить"
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Купить", callbackData),
		),
	)
}
