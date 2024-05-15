package model

import "time"

type User struct {
	UserID       string    `json:"user_id"`
	NIP          string    `json:"nip"`
	Name         string    `json:"name"`
	PasswordHash []byte    `json:"password_hash"`
	Role         string    `json:"role"`
	CardImage    string    `json:"card_image"`
	CreatedAt    time.Time `json:"created_at"`
}
