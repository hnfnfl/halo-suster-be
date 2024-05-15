package model

import "time"

type Medical struct {
	UniqueID       string    `json:"unique_id"`
	IdentityNumber string    `json:"identity_number"`
	CreatorID      string    `json:"creator_id"`
	Symptoms       string    `json:"symptoms"`
	Medication     string    `json:"medication"`
	CreatedAt      time.Time `json:"created_at"`
}
