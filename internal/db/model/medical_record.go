package model

import "time"

type MedicalRecord struct {
	UniqueID       string    `db:"unique_id"`
	IdentityNumber string    `db:"identity_number"`
	CreatorID      string    `db:"creator_id"`
	Symptoms       string    `db:"symptoms"`
	Medication     string    `db:"medication"`
	CreatedAt      time.Time `db:"created_at"`
}
