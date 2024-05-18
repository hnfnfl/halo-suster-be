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
		errs.NewBadRequestError("JSON binding error", err).Send(ctx)
		return
	}

	if err := request.Validate(); err != nil {
		errs.NewValidationError("input validation error", err).Send(ctx)
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

func (mrh *MedicalRecordHandler) GetMedicalRecord(ctx *gin.Context) {
	var param dto.ReqParamGetMedicalRecord
	// var limit int
	param.IdentityNumber = ctx.Query("identityDetail.identityNumber")
	if ctx.Query("limit") != "" {
		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			errs.NewBadRequestError("error convert param limit", err).Send(ctx)
			return
		}
		param.Limit = limit
	} else {
		param.Limit = 5
	}
	if ctx.Query("offset") != "" {
		offset, err := strconv.Atoi(ctx.Query("offset"))
		if err != nil {
			errs.NewBadRequestError("error convert param offset", err).Send(ctx)
			return
		}
		param.Offset = offset
	} else {
		param.Offset = 0
	}
	// offset, err := strconv.Atoi(ctx.Query("offset"))
	// if err != nil {
	// 	errs.NewInternalError("error convert param offset", err).Send(ctx)
	// 	return
	// }
	// param.Offset = offset
	param.UserId = ctx.Query("createdBy.userId")
	param.Nip = ctx.Query("createdBy.nip")
	param.CreatedAt = dto.Sort(ctx.Query("created_at"))
	mrh.service.GetAllMedicalRecord(param).Send(ctx)
}
