// message_service.go — сервис для обработки исходящих сообщений для пользователя
// без разницы куда послать сообщение, телеграм бот, сайт, приложение
package service

// MessageService - сервис для генерации сообщений
type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) GetWelcomeMessage() string {
	return "Добро пожаловать в Salle Parfume! Я помогу тебе найти идеальный парфюм."
}

func (s *MessageService) GetAboutMessage() string {
	return "Мы - Salle Parfume, ваш проводник в мир ароматов."
}
