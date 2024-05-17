package handler

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"halo-suster/internal/pkg/middleware"
	"halo-suster/internal/pkg/service"
	"halo-suster/internal/pkg/util"
	"strconv"

	// "net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service   *service.Service
	validator *validator.Validate
}

func NewUserHandler(s *service.Service, validator *validator.Validate) *UserHandler {
	return &UserHandler{s, validator}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	in := dto.RegisterRequest{}
	msg, err := util.JsonBinding(ctx, &in)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	// validate Request
	if err := in.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	data := model.User{
		NIP:  strconv.Itoa(in.NIP),
		Name: in.Name,
	}

	var passHash []byte
	if in.Password != "" {
		var err error
		passHash, err = middleware.PasswordHash(in.Password, h.service.Config().Salt)
		if err != nil {
			errs.NewInternalError("hashing error", err).Send(ctx)
			return
		}
	}

	if strconv.Itoa(in.NIP)[:3] == "615" {
		// validate as IT
		data.UserID = util.UuidGenerator(util.ITPrefix)
		data.Role = "it"
		data.PasswordHash = passHash
	} else if strconv.Itoa(in.NIP)[:3] == "303" {
		// validate as Nurse
		data.UserID = util.UuidGenerator(util.NursePrefix)
		data.Role = "nurse"
		data.CardImage = *in.CardImage
	}

	h.service.RegisterUser(data).Send(ctx)
}

func (h *UserHandler) ITLogin(ctx *gin.Context) {
	body := dto.LoginRequest{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		errs.NewInternalError("JSON binding error", err).Send(ctx)
		return
	}

	err := h.validator.Struct(body)
	if err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	data := model.User{}

	if strconv.Itoa(body.NIP)[:3] == "615" {
		data.NIP = strconv.Itoa(body.NIP)
		data.Role = "it"
	} else {
		errs.NewNotFoundError("user is not from IT (nip not starts with 615)")
		return
	}

	h.service.LoginUser(data, body).Send(ctx)

}

func (h *UserHandler) NurseLogin(ctx *gin.Context) {
	body := dto.LoginRequest{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		errs.NewInternalError("JSON binding error", err).Send(ctx)
		return
	}

	err := h.validator.Struct(body)
	if err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	data := model.User{}
	var passHash []byte
	if body.Password != "" {
		var err error
		passHash, err = middleware.PasswordHash(body.Password, h.service.Config().Salt)
		if err != nil {
			errs.NewInternalError("hashing error", err).Send(ctx)
			return
		}
	}

	if strconv.Itoa(body.NIP)[:3] == "615" {
		data.NIP = strconv.Itoa(body.NIP)
		data.Role = "it"
		data.PasswordHash = passHash
	} else {
		errs.NewNotFoundError("user is not from IT (nip not starts with 615)")
		return
	}

	h.service.LoginUser(data, body).Send(ctx)

}
