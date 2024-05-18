package service

import (
	"fmt"
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"strings"
	"time"
)

func (s *Service) CreateMedicalRecord(medicalRecord model.MedicalRecord) errs.Response {
	var err error
	stmt := "INSERT INTO medical_records (unique_id, identity_number, creator_id, symptoms, medication, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	if _, err = s.DB().Exec(stmt, medicalRecord.UniqueID, medicalRecord.IdentityNumber, medicalRecord.CreatorID, medicalRecord.Symptoms, medicalRecord.Medication, time.Now()); err != nil {
		return errs.NewInternalError("insert medical record error", err)
	}
	return errs.Response{
		Message: "Medical record successfully added",
		Data: dto.ResponseCreateMedicalRecord{
			IdentityNumber: medicalRecord.IdentityNumber,
		},
	}
}

func (s *Service) GetAllMedicalRecord(param dto.ReqParamGetMedicalRecord) errs.Response {
	var err error
	var query strings.Builder
	var createdAt time.Time
	var birthDate time.Time

	query.WriteString("SELECT p.identity_number , p.phone_number , p.name , p.birth_date , p.gender , p.card_image , mr.symptoms , mr.medication , mr.created_at , u.nip , u.name , u.user_id  FROM medical_records mr join patients p on mr.identity_number = p.identity_number join users u ON mr.creator_id = u.user_id WHERE 1=1 ")
	// if param.UniqueID != "" {
	// 	query.WriteString("AND unique_id = " + param.UniqueID)
	// }
	if param.IdentityNumber != "" {
		query.WriteString("AND p.identity_number = " + param.IdentityNumber)
	}
	if param.UserId != "" {
		query.WriteString("AND u.user_id = " + param.UserId)
	}
	if param.Nip != "" {
		query.WriteString("AND u.nip = " + param.Nip)
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
		medicalRecord.Detail.BirthDate = birthDate.Format(time.RFC3339)
		medicalRecord.CreatedAt = createdAt.Format(time.RFC3339)
		medicalRecords = append(medicalRecords, medicalRecord)
	}

	return errs.Response{
		Message: "successfully get medical records",
		Data:    medicalRecords,
	}
}
