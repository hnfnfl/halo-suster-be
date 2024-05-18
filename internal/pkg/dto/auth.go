package dto

import (
	"strconv"
	"time"

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

	// 	- the fourth digit, if it's male, fill it with `1`, else `2`
	genderUser, err := strconv.Atoi(nipStr[3:4])
	if err != nil {
		return validation.NewError("nip", "NIP is not valid, the fourth digit, if it's male, fill it with `1`, else `2`")
	}
	if genderUser < 1 || genderUser > 2 {
		return validation.NewError("nip", "NIP is not valid, the fourth digit, if it's male, fill it with `1`, else `2`")
	}

	// 	- the fifth and eigth digit, fill it with a year, starts from `2000` till current year
	yearUser, err := strconv.Atoi(nipStr[4:8])
	currentYear := time.Now().Year()
	if err != nil {
		return validation.NewError("nip", "NIP is not valid, the fifth and eigth digit, fill it with a year, starts from `2000` till current year")
	}
	if yearUser < 2000 || yearUser > currentYear {
		return validation.NewError("nip", "NIP is not valid, the fifth and eigth digit, fill it with a year, starts from `2000` till current year")
	}

	// 	- the ninth and tenth, fill it with month, starts from `01` till `12`
	monthUser, err := strconv.Atoi(nipStr[8:10])
	if err != nil {
		return validation.NewError("nip", "NIP is not valid, the ninth and tenth, fill it with month, starts from `01` till `12`")
	}

	if monthUser < 1 || monthUser > 12 {
		return validation.NewError("nip", "NIP is not valid, the ninth and tenth, fill it with month, starts from `01` till `12`")
	}

	if err := validation.ValidateStruct(&r,
		validation.Field(&r.NIP, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(5, 50)),
	); err != nil {
		return err
	}

	// - first until third digit, should start with `615`
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

func (r LoginRequest) Validate() error {
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
		validation.Field(&r.Password, validation.Required, validation.Length(5, 33)),
	); err != nil {
		return err
	}
	return nil
}
