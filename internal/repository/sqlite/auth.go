// auth.go - реализация интерфейса Authorization для SQLite
package sqlite

import (
	"database/sql"
	"fmt"
	"salle_parfume/internal/domain"
	"salle_parfume/internal/repository"
)

type AuthSqlite struct {
	db *sql.DB
}

func NewAuthSqlite(db *sql.DB) repository.Authorization {
	// Инициализация таблицы при создании репозитория
	// В продакшене лучше использовать миграции, но для старта это ок
	if err := createUsersTable(db); err != nil {
		// Логируем ошибку, но не паникуем, хотя можно и вернуть error из конструктора
		fmt.Printf("Error creating users table: %v\n", err)
	}
	return &AuthSqlite{db: db}
}

func createUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL UNIQUE,
		username TEXT,
		first_name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	return err
}

func (r *AuthSqlite) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (chat_id, username, first_name) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, user.ChatID, user.Username, user.FirstName)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *AuthSqlite) GetUserByChatID(chatID int64) (*domain.User, error) {
	query := `SELECT id, chat_id, username, first_name, created_at FROM users WHERE chat_id = ?`
	var user domain.User
	err := r.db.QueryRow(query, chatID).Scan(&user.ID, &user.ChatID, &user.Username, &user.FirstName, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
