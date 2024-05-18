package handler

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"halo-suster/internal/pkg/middleware"
	"halo-suster/internal/pkg/service"
	"halo-suster/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type NurseHandler struct {
	service   *service.Service
	validator *validator.Validate
}

func NewNurseHandler(s *service.Service, validator *validator.Validate) *NurseHandler {
	return &NurseHandler{s, validator}
}

func (nh *NurseHandler) UpdateNurse(ctx *gin.Context) {
	request := dto.RequestUpdateNurse{}
	msg, err := util.JsonBinding(ctx, &request)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	// validate Request
	if err := request.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	userID := ctx.Param("userId")
	role := ctx.Value("userRole").(string)

	if role == "it" {
		nh.service.UpdateNurse(request).Send(ctx)
	} else {
		errs.NewUnauthorizedError("user is not authorized").Send(ctx)
		return
	}
}

func (nh *NurseHandler) DeleteNurse(ctx *gin.Context) {

	userID := ctx.Param("userId")
	role := ctx.Value("userRole").(string)

	if role == "it" {
		nh.service.DeleteNurse(userID).Send(ctx)
	} else {
		errs.NewUnauthorizedError("user is not authorized").Send(ctx)
		return
	}
}

func (nh *NurseHandler) AccessNurse(ctx *gin.Context) {
	request := dto.RequestCreateAccessNurse{}
	msg, err := util.JsonBinding(ctx, &request)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	// validate Request
	if err := request.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	data := model.User{}
	var passHash []byte
	var errHash error

	passHash, errHash = middleware.PasswordHash(request.Password, nh.service.Config().Salt)
	if errHash != nil {
		errs.NewInternalError("hashing error", errHash).Send(ctx)
		return
	}
	data.PasswordHash = passHash
	data.UserID = ctx.Param("userId")

	role := ctx.Value("userRole").(string)

	if role == "it" {
		nh.service.AccessNurse(data).Send(ctx)
	} else {
		errs.NewUnauthorizedError("user is not authorized").Send(ctx)
		return
	}
}
