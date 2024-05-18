package model

import "time"

var NIPUsers = map[string]bool{
	"605": true,
	"303": true,
}

type User struct {
	UserID       string    `db:"user_id"`
	NIP          string    `db:"nip"`
	Name         string    `db:"name"`
	PasswordHash []byte    `db:"password_hash"`
	Role         string    `db:"role"`
	CardImage    string    `db:"card_image"`
	CreatedAt    time.Time `db:"created_at"`
}
