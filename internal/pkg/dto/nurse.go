package dto

import (
	"errors"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RequestUpdateNurse struct {
	NIP  uint64 `json:"nip"`
	Name string `json:"name" validate:"required,min=5,max=50"`
}

type RequestCreateAccessNurse struct {
	Password string `json:"password" validate:"required,min=5,max=33"`
}

func (r RequestCreateAccessNurse) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Password, validation.Required, validation.Length(5, 33)),
	); err != nil {
		return err
	} else {
		return nil
	}
}

func regexMatch(re *regexp.Regexp) validation.RuleFunc {
	return func(value interface{}) error {
		s := fmt.Sprintf("%v", value)
		if !re.MatchString(s) {
			return errors.New("string doesn't match the rules")
		}
		return nil
	}
}

func (r RequestUpdateNurse) Validate() error {
	nipRegex := regexp.MustCompile("^(303|615)[12][2-9][0-9]{3}(0[1-9]|1[0-2])[0-9]{3,5}$")

	return validation.ValidateStruct(&r,
		validation.Field(&r.NIP, validation.Required, validation.By(regexMatch(nipRegex))),
		validation.Field(&r.Name, validation.Required, validation.Length(5, 50)),
	)
}
