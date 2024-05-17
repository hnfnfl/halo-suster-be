package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RequestCreateAccessNurse struct {
	Password string `json:"password" validate:"required,min=5,max=33"`
}

func (r RequestCreateAccessNurse) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Password, validation.Required, validation.Length(5, 33)),
	); err != nil {
		return err
	}

	return validation.NewError("password", "Password is required")
}
