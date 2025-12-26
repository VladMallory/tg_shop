// product.go - Реализация интерфейса ProductRepository для базы данных SQLite.
// Здесь мы пишем конкретные SQL-запросы.
package sqlite

import (
	"database/sql"
	"fmt"
	"salle_parfume/internal/domain"
	"salle_parfume/internal/repository"
)

// ProductSqlite - структура, хранящая подключение к БД.
// Она реализует интерфейс repository.ProductRepository.
type ProductSqlite struct {
	db *sql.DB
}

// NewProductSqlite - создает новый экземпляр репозитория товаров.
func NewProductSqlite(db *sql.DB) repository.ProductRepository {
	// Сразу при запуске проверяем, есть ли таблица products.
	// Если нет - создаем её.
	if err := createProductsTable(db); err != nil {
		// Если не удалось создать таблицу - пишем в консоль, но не роняем программу (хотя можно и запаниковать)
		fmt.Printf("Error creating products table: %v\n", err)
	}
	return &ProductSqlite{db: db}
}

// createProductsTable - SQL запрос для создания таблицы товаров
func createProductsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT,         -- Тип (female, male, unisex)
		name TEXT,         -- Название
		description TEXT,  -- Описание
		price REAL,        -- Цена (REAL это float в sqlite)
		image_id TEXT      -- ID картинки в телеграм
	);
	`
	_, err := db.Exec(query)
	return err
}

// CreateProduct - Добавляет товар в базу данных
func (r *ProductSqlite) CreateProduct(product *domain.Product) error {
	// Используем подготовленные выражения (?) для защиты от SQL-инъекций
	query := `INSERT INTO products (type, name, description, price, image_id) VALUES (?, ?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, product.Type, product.Name, product.Description, product.Price, product.ImageID)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetAllProducts - Получает список всех товаров из базы
func (r *ProductSqlite) GetAllProducts() ([]domain.Product, error) {
	query := `SELECT id, type, name, description, price, image_id FROM products`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Обязательно закрываем rows, чтобы не текли соединения

	var products []domain.Product
	
	// Бежим по строкам результата
	for rows.Next() {
		var p domain.Product
		// Сканируем данные из строки в структуру
		if err := rows.Scan(&p.ID, &p.Type, &p.Name, &p.Description, &p.Price, &p.ImageID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}