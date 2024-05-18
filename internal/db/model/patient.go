package model

import "time"

type PatientGender string

const (
	Male   PatientGender = "male"
	Female PatientGender = "female"
)

var Genders = []string{"male", "female"}

type Patient struct {
	IdentityNumber string        `db:"user_id"`
	Name           string        `db:"name"`
	BirthDate      string        `db:"birth_date"`
	PhoneNumber    string        `db:"phone_number"`
	CardImage      string        `db:"card_image"`
	Gender         PatientGender `db:"gender"`
	CreatedAt      time.Time     `db:"created_at"`
}
