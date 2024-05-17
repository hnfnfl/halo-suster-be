package dto

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RegisterRequest struct {
	NIP       int     `json:"nip"`
	Name      string  `json:"name"`
	Password  string  `json:"password"`
	CardImage *string `json:"identityCardScanImg"`
}

type LoginRequest struct {
	NIP      int    `json:"nip" validate:"required,min=13,max=13"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type AuthResponse struct {
	UserId      string `json:"userId"`
	NIP         string `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken,omitempty"`
}

func (r RegisterRequest) Validate() error {
	if err := validation.Validate(r.NIP, validation.Required); err != nil {
		return validation.NewError("nip", "NIP is required")
	}

	// validate NIP
	nipStr := strconv.Itoa(r.NIP)
	if len(nipStr) != 13 {
		return validation.NewError("nip", "NIP is not valid")
	}

	if err := validation.ValidateStruct(&r,
		validation.Field(&r.NIP, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(5, 50)),
	); err != nil {
		return err
	}
	if nipStr[:3] == "615" {
		// validate as IT
		return validation.ValidateStruct(&r, validation.Field(&r.Password, validation.Required, validation.Length(5, 33)))
	} else if nipStr[:3] == "303" {
		// validate as Nurse
		return validation.ValidateStruct(&r, validation.Field(&r.CardImage, validation.Required))
	} else {
		return validation.NewError("nip", "NIP is not valid")
	}
}
