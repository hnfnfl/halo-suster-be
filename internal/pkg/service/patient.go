package service

import (
	"fmt"
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"net/http"
	"strings"
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

func (s *Service) Get(param dto.ReqParamGetPatient) errs.Response {
	var err error
	var query strings.Builder
	var createdAt time.Time
	var birthDate time.Time

	query.WriteString("SELECT identity_number, name, birth_date, phone_number, gender, created_at FROM patients WHERE 1=1 ")
	if param.IdentityNumber != "" {
		query.WriteString("AND identity_number = " + param.IdentityNumber)
	}
	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND name LIKE '%%%s%%' ", strings.ToLower(param.Name)))
	}
	if param.PhoneNumber != "" {
		query.WriteString("AND phone_number = +" + param.PhoneNumber)
	}

	if param.CreatedAt == "asc" {
		query.WriteString("ORDER BY created_at ASC ")
	} else {
		query.WriteString("ORDER BY created_at DESC ")
	}

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))
	rows, err := s.DB().Query(query.String())
	if err != nil {
		return errs.NewInternalError("query error", err)
	}
	defer rows.Close()

	results := []dto.ResponseGetPatient{}

	for rows.Next() {
		result := dto.ResponseGetPatient{}
		err := rows.Scan(
			&result.IdentityNumber,
			&result.Name,
			&birthDate,
			&result.PhoneNumber,
			&result.Gender,
			&createdAt,
		)
		result.BirthDate = birthDate.Format(time.RFC3339)
		result.CreatedAt = createdAt.Format(time.RFC3339)
		if err != nil {
			return errs.NewInternalError("error scan query db", err)
		}
		results = append(results, result)
	}
	return errs.Response{
		Code:    http.StatusOK,
		Message: "successfully get patients",
		Data:    results,
	}
}
