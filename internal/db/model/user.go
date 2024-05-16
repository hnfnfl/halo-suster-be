package model

import "time"

type User struct {
	UserID       string    `db:"user_id"`
	NIP          string    `db:"nip"`
	Name         string    `db:"name"`
	PasswordHash []byte    `db:"password_hash"`
	Role         string    `db:"role"`
	CardImage    string    `db:"card_image"`
	CreatedAt    time.Time `db:"created_at"`
}
