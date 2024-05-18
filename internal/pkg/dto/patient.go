package dto

import (
	"halo-suster/internal/db/model"
	"regexp"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func isISO8601Date(value interface{}) error {
	str, _ := value.(string)
	_, err := time.Parse("2006-01-02", str)
	if err != nil {
		return validation.NewError("validation_iso8601_date", "must be a valid ISO 8601 date (yyyy-mm-dd)")
	}
	return nil
}
func isISO8601Datetime(value interface{}) error {
	str, _ := value.(string)
	_, err := time.Parse("2006-01-02T15:04:05Z07:00", str)
	if err != nil {
		return validation.NewError("validation_iso8601_datetime", "must be a valid ISO 8601 date (yyyy-mm-dd)")
	}
	return nil
}

var phoneNumberValidationRule = validation.NewStringRule(func(s string) bool {
	return strings.HasPrefix(s, "+")
}, "phone number must start with international calling code")

var imgUrlValidationRule = validation.NewStringRule(func(s string) bool {
	match, _ := regexp.MatchString(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/{1}[A-z0-9_\-\:x\=\(\)]+)*(\.(jpg|jpeg|png))?$`, s)
	return match
}, "image url is not valid")

type RequestCreatePatient struct {
	IdentityNumber      *int   `json:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber"`
	Name                string `json:"name"`
	BirthDate           string `json:"birthDate"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}

type ResponseCreatePatient struct {
	Name string `json:"name"`
}

func (r RequestCreatePatient) Validate() error {
	var identityNumberStr string
	if r.IdentityNumber != nil {
		identityNumberStr = strconv.Itoa(*r.IdentityNumber)
	}
	if len(identityNumberStr) != 16 {
		return validation.NewError("identityNumber", "identityNumber is not valid")
	}

	// if(len() != 16) {
	// 	return validation.NewError("identityNumber", "Identity Number must be 16 digits")
	// }
	// err := validation.Validate(result, validation.Required, validation.Length(16, 16))
	// if err != nil {
	// 	return validation.NewError("identityNumber", "Identity Number must be 16 digits")
	// }
	return validation.ValidateStruct(&r,
		validation.Field(&r.IdentityNumber, validation.Required),
		validation.Field(&r.PhoneNumber, validation.Required, phoneNumberValidationRule, validation.Length(10, 16)),
		validation.Field(&r.Name, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.BirthDate, validation.Required, validation.By(isISO8601Datetime)),
		validation.Field(&r.Gender, validation.Required, validation.In(model.Genders)),
		validation.Field(&r.IdentityCardScanImg, validation.Required, imgUrlValidationRule),
	)
}

type ReqParamGetPatient struct {
	IdentityNumber string `json:"identityNumber"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	Name           string `json:"name"`
	PhoneNumber    string `json:"phoneNumber"`
	CreatedAt      Sort   `json:"createdAt"`
}

type ResponseGetPatient struct {
	IdentityNumber int    `json:"identityNumber"`
	PhoneNumber    string `json:"phoneNumber"`
	Name           string `json:"name"`
	BirthDate      string `json:"birthDate"`
	Gender         string `json:"gender"`
	CreatedAt      string `json:"createdAt"`
}
