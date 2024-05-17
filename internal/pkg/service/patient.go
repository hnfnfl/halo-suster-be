package service

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"time"
)

func (s *Service) Create(patient model.Patient) errs.Response {
	var err error

	tx, err := s.DB().Begin()
	if err != nil {
		return errs.NewInternalError("transaction error", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}()

	stmt := "INSERT INTO patients (identity_number,name,birth_date, phone_number,gender,card_image,created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if _, err = tx.Exec(stmt, patient.IdentityNumber, patient.Name, patient.BirthDate, patient.PhoneNumber, patient.Gender, patient.CardImage, time.Now()); err != nil {
		return errs.NewInternalError("insert error", err)
	}
	return errs.Response{
		Message: "Medical patient successfully added",
		Data: dto.ResponseCreatePatient{
			Name: patient.Name,
		},
	}
}