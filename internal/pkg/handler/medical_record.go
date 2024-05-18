package handler

import (
	"fmt"
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"halo-suster/internal/pkg/service"
	"halo-suster/internal/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MedicalRecordHandler struct {
	service   *service.Service
	validator *validator.Validate
}

func NewMedicalRecordHandler(s *service.Service, validator *validator.Validate) *MedicalRecordHandler {
	return &MedicalRecordHandler{s, validator}
}

func (mrh *MedicalRecordHandler) CreateMedicalRecord(ctx *gin.Context) {
	request := dto.RequestCreateMedicalRecord{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errs.NewInternalError("JSON binding error", err).Send(ctx)
		return
	}

	//validate request input
	if err := request.Validate(); err != nil {
		errs.NewValidationError("input validation error", err).Send(ctx)
		return
	}

	//get user role
	userId := ctx.Value("userID").(string)
	fmt.Print("ini usernip : ", userId)

	data := model.MedicalRecord{
		UniqueID:       util.UuidGenerator(""),
		IdentityNumber: strconv.Itoa(*request.IdentityNumber),
		CreatorID:      userId,
		Symptoms:       request.Symptoms,
		Medication:     request.Medications,
	}
	mrh.service.CreateMedicalRecord(data).Send(ctx)
}
