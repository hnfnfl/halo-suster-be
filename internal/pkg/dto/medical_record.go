package dto

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RequestCreateMedicalRecord struct {
	IdentityNumber *int   `json:"identityNumber"`
	Symptoms       string `json:"symptoms"`
	Medications    string `json:"medications"`
}

func (r RequestCreateMedicalRecord) Validate() error {
	result := strconv.Itoa(*r.IdentityNumber)
	err := validation.Validate(result, validation.Required, validation.Length(16, 16))
	if err != nil {
		return validation.NewError("identityNumber", "Identity Number must be 16 digits")
	}
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Symptoms, validation.Required, validation.Length(1, 2000)),
		validation.Field(&r.Medications, validation.Required, validation.Length(1, 2000)),
	); err != nil {
		return err
	} else {
		return nil
	}
}

type ResponseCreateMedicalRecord struct {
	IdentityNumber string `json:"identityNumber"`
}
