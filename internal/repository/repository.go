// repository.go - Интерфейсы для работы с хранилищем данных.
// Repository Layer отвечает за сохранение и загрузку данных.
// Интерфейсы позволяют нам менять базу данных (например, с SQLite на Postgres) не меняя остальной код.
package repository

import "salle_parfume/internal/domain"

// Authorization - Контракт для работы с пользователями.
type Authorization interface {
	CreateUser(user *domain.User) error                 // Сохранить нового пользователя
	GetUserByChatID(chatID int64) (*domain.User, error) // Найти пользователя по ID чата
}

// ProductRepository - Контракт для работы с товарами (Духами).
// Мы описываем ЧТО мы хотим делать, но не КАК.
type ProductRepository interface {
	CreateProduct(product *domain.Product) error // Сохранить товар
	GetAllProducts() ([]domain.Product, error)   // Получить список всех товаров
}

// Repository - Главная структура, которая объединяет все наши репозитории.
// Это удобно, чтобы передавать один объект `Repository` в Handler, вместо кучи мелких.
type Repository struct {
	Authorization
	ProductRepository
}

// NewRepository - Конструктор. Собирает отдельные реализации в одну коробку.
// Принимает:
// auth - реализацию работы с юзерами
// prod - реализацию работы с товарами
func NewRepository(auth Authorization, prod ProductRepository) *Repository {
	return &Repository{
		Authorization:     auth,
		ProductRepository: prod,
	}
}
