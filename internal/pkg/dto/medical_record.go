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
	var identityNumberStr string
	if r.IdentityNumber != nil {
		identityNumberStr = strconv.Itoa(*r.IdentityNumber)
	}
	if len(identityNumberStr) != 16 {
		return validation.NewError("identityNumber", "identityNumber is not valid")
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

type ReqParamGetMedicalRecord struct {
	IdentityNumber string
	UserId         string
	Nip            string
	Limit          int `json:"limit"`
	Offset         int `json:"offset"`
	CreatedAt      Sort
}

type ResponseGetMedicalRecord struct {
	Detail      IdentityDetail `json:"identityDetail"`
	Symptoms    string         `json:"symptoms"`
	Medications string         `json:"medications"`
	CreatedAt   string         `json:"createdAt"`
	CreatedBy   CreatedBy      `json:"createdBy"`
}

type IdentityDetail struct {
	IdentityNumber      int    `json:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber"`
	Name                string `json:"name"`
	BirthDate           string `json:"birthDate"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}

type CreatedBy struct {
	Nip    int    `json:"nip"`
	Name   string `json:"name"`
	UserId string `json:"userId"`
}
