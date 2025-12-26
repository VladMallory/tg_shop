// sqlite.go - инициализация подключения к SQLite
package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Драйвер для SQLite
)

type Config struct {
	DriverName string
	Path       string
}

func NewSqliteDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.DriverName, cfg.Path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
