package service

import (
	"fmt"
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"time"
)

func (s *Service) CreateMedicalRecord(medicalRecord model.MedicalRecord) errs.Response {
	var err error
	fmt.Println("sudah disini")
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
