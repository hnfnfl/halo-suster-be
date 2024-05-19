package service

import (
	"errors"
	"fmt"
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"net/http"
	"strings"
	"time"
)

func (s *Service) CreateMedicalRecord(medicalRecord model.MedicalRecord) errs.Response {
	// var err error
	if !s.isIdentityNumberExist(medicalRecord.IdentityNumber) {
		return errs.NewNotFoundError("identity number not found", errors.New("identity number not found"))
	}
	stmt := "INSERT INTO medical_records (unique_id, identity_number, creator_id, symptoms, medication, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, errQuery := s.DB().Exec(stmt, medicalRecord.UniqueID, medicalRecord.IdentityNumber, medicalRecord.CreatorID, medicalRecord.Symptoms, medicalRecord.Medication, time.Now())
	if errQuery != nil {
		return errs.NewInternalError("insert medical record error", errQuery)

	}
	// rowAffected, _ := rows.RowsAffected()
	// if rowAffected == 0 {
	// 	return errs.NewBadRequestError("identity number not found", nil)
	// }
	return errs.Response{
		Code:    http.StatusCreated,
		Message: "Medical record successfully added",
		Data: dto.ResponseCreateMedicalRecord{
			IdentityNumber: medicalRecord.IdentityNumber,
		},
	}
}

func (s *Service) isIdentityNumberExist(identityNumber string) bool {
	query := "SELECT identity_number FROM patients WHERE identity_number = $1"
	var identityNumberExist string
	err := s.DB().QueryRow(query, identityNumber).Scan(&identityNumberExist)
	if err != nil {
		return false
	}
	return identityNumberExist == identityNumber
	// if err :=  err != nil {
	// 	return identityNumberExist != *identityNumber
	// }
	// return identityNumberExist == *identityNumber
}

func (s *Service) GetAllMedicalRecord(param dto.ReqParamGetMedicalRecord) errs.Response {
	var err error
	var query strings.Builder
	var createdAt time.Time
	var birthDate time.Time

	query.WriteString("SELECT p.identity_number , p.phone_number , p.name , p.birth_date , p.gender , p.card_image , mr.symptoms , mr.medication , mr.created_at , u.nip , u.name , u.user_id  FROM medical_records mr join patients p on mr.identity_number = p.identity_number join users u ON mr.creator_id = u.nip WHERE 1=1 ")
	// if param.UniqueID != "" {
	// 	query.WriteString("AND unique_id = " + param.UniqueID)
	// }
	if param.IdentityNumber != "" {
		query.WriteString("AND p.identity_number = " + param.IdentityNumber)
	}
	if param.UserId != "" {
		query.WriteString(fmt.Sprintf("AND u.user_id = '%s'", param.UserId))
	}
	if param.Nip != "" {
		query.WriteString(fmt.Sprintf("AND u.nip = '%s'", param.Nip))
	}

	if param.CreatedAt == "asc" {
		query.WriteString(" ORDER BY created_at ASC ")
	} else {
		query.WriteString(" ORDER BY created_at DESC ")
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

	var medicalRecords []dto.ResponseGetMedicalRecord
	for rows.Next() {
		var medicalRecord dto.ResponseGetMedicalRecord
		if err = rows.Scan(
			&medicalRecord.Detail.IdentityNumber,
			&medicalRecord.Detail.PhoneNumber,
			&medicalRecord.Detail.Name,
			&birthDate,
			&medicalRecord.Detail.Gender,
			&medicalRecord.Detail.IdentityCardScanImg,
			&medicalRecord.Symptoms,
			&medicalRecord.Medications,
			&createdAt,
			&medicalRecord.CreatedBy.Nip,
			&medicalRecord.CreatedBy.Name,
			&medicalRecord.CreatedBy.UserId,
		); err != nil {
			return errs.NewInternalError("scan error", err)
		}
		medicalRecord.Detail.BirthDate = birthDate.Format(time.RFC3339Nano)
		medicalRecord.CreatedAt = createdAt.Format(time.RFC3339Nano)
		medicalRecords = append(medicalRecords, medicalRecord)
	}

	return errs.Response{
		Message: "successfully get medical records",
		Data:    medicalRecords,
	}
}
