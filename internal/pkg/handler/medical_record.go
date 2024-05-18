package handler

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"halo-suster/internal/pkg/service"
	"halo-suster/internal/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MedicalRecordHandler struct {
	service *service.Service
}

func NewMedicalRecordHandler(s *service.Service) *MedicalRecordHandler {
	return &MedicalRecordHandler{s}
}

func (mrh *MedicalRecordHandler) CreateMedicalRecord(ctx *gin.Context) {
	request := dto.RequestCreateMedicalRecord{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errs.NewInternalError("JSON binding error", err).Send(ctx)
		return
	}

	//get user role
	userId := ctx.Value("userID").(string)

	data := model.MedicalRecord{
		UniqueID:       util.UuidGenerator(""),
		IdentityNumber: strconv.Itoa(*request.IdentityNumber),
		CreatorID:      userId,
		Symptoms:       request.Symptoms,
		Medication:     request.Medications,
	}
	mrh.service.CreateMedicalRecord(data).Send(ctx)
}
