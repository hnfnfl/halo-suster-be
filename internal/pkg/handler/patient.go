package handler

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"halo-suster/internal/pkg/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	*service.Service
}

func NewPatientHandler(s *service.Service) *PatientHandler {
	return &PatientHandler{s}
}

func (ph *PatientHandler) CreatePatient(ctx *gin.Context) {
	request := dto.RequestCreatePatient{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errs.NewInternalError("JSON binding error", err).Send(ctx)
		return
	}

	//validate request input
	if err := request.Validate(); err != nil {
		errs.NewValidationError("input validation error", err).Send(ctx)
		return
	}

	data := model.Patient{
		IdentityNumber: strconv.Itoa(*request.IdentityNumber),
		Name:           request.Name,
		BirthDate:      request.BirthDate,
		PhoneNumber:    request.PhoneNumber,
		CardImage:      request.IdentityCardScanImg,
		Gender:         model.PatientGender(request.Gender),
		CreatedAt:      time.Now(),
	}

	ph.Create(data).Send(ctx)
}

func (ph *PatientHandler) GetPatient(ctx *gin.Context) {
	var param dto.ReqParamGetPatient
	param.IdentityNumber = ctx.Query("identity_number")
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		errs.NewInternalError("error convert param limit", err).Send(ctx)
		return
	}
	param.Limit = limit
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		errs.NewInternalError("error convert param offset", err).Send(ctx)
		return
	}
	param.Offset = offset
	param.Name = ctx.Query("name")
	param.PhoneNumber = ctx.Query("phone_number")
	param.CreatedAt = dto.Sort(ctx.Query("created_at"))
	ph.Get(param).Send(ctx)
}
