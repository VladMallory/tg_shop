// user.go чистая модель данных. То из чего состоит сущность юзеров
package domain

import "time"

// User - основная сущность пользователя
type User struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	CreatedAt time.Time `json:"created_at"`
}
