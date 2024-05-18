package handler

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"halo-suster/internal/pkg/middleware"
	"halo-suster/internal/pkg/service"
	"halo-suster/internal/pkg/util"
	"strconv"
	"strings"

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
	body := dto.RegisterRequest{}
	msg, err := util.JsonBinding(ctx, &body)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	// validate Request
	if err := body.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	data := model.User{
		NIP:  strconv.Itoa(body.NIP),
		Name: body.Name,
	}

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
		// validate as IT
		data.UserID = util.UuidGenerator(util.ITPrefix)
		data.Role = "it"
		data.PasswordHash = passHash
	} else if strconv.Itoa(body.NIP)[:3] == "303" {
		// validate as Nurse
		data.UserID = util.UuidGenerator(util.NursePrefix)
		data.Role = "nurse"
		data.CardImage = *body.CardImage
	}
	role := extractRole(ctx.FullPath())

	switch role {
	case "it":
		h.service.RegisterUser(data).Send(ctx)
	case "nurse":
		role := ctx.Value("userRole").(string)
		if role == "it" {
			h.service.RegisterUser(data).Send(ctx)
		} else {
			errs.NewUnauthorizedError("user is not authorized").Send(ctx)
			return
		}
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	body := dto.LoginRequest{}
	msg, err := util.JsonBinding(ctx, &body)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	// if err := h.validator.Struct(body); err != nil {
	// 	errs.NewValidationError("Request validation error", err).Send(ctx)
	// 	return
	// }

	data := model.User{
		NIP: strconv.Itoa(body.NIP),
	}
	role := extractRole(ctx.FullPath())
	if body.Password != "" {
		data.PasswordHash = []byte(body.Password + h.service.Config().Salt)
	}

	switch role {
	case "it":
		// validate Request
		if err := body.Validate(); err != nil {
			errs.NewValidationError("Request validation error", err).Send(ctx)
			return
		}
		if strconv.Itoa(body.NIP)[:3] == "615" {
			data.Role = "it"
		} else {
			errs.NewNotFoundError("user is not from IT (nip not starts with 615)", errs.ErrUserNotFound).Send(ctx)
			return
		}
	case "nurse":
		// validate Request
		if err := body.Validate(); err != nil {
			errs.NewValidationError("Request validation error", err).Send(ctx)
			return
		}
		if strconv.Itoa(body.NIP)[:3] == "303" {
			data.Role = "nurse"
		} else {
			errs.NewNotFoundError("user is not from IT (nip not starts with 615)", errs.ErrUserNotFound).Send(ctx)
			return
		}
	}

	h.service.LoginUser(data).Send(ctx)
}

func extractRole(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		return parts[3]
	}
	return ""
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	queryParams := ctx.Request.URL.Query()
	var param dto.ReqParamUserGet

	param.UserID = queryParams.Get("userId")
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		errs.NewBadRequestError("param limit should be a number", errs.ErrBadParam).Send(ctx)
		return
	}
	param.Limit = limit
	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil {
		errs.NewBadRequestError("param offset should be a number", errs.ErrBadParam).Send(ctx)
		return
	}
	param.Offset = offset

	param.Name = queryParams.Get("name")
	if queryParams.Get("nip") != "" {
		_, err := strconv.Atoi(queryParams.Get("nip"))
		if err != nil {
			errs.NewBadRequestError("param nip should be a number", errs.ErrBadParam).Send(ctx)
			return
		} else {
			param.NIP = queryParams.Get("nip")
		}
	}
	param.Role = dto.Role(queryParams.Get("role"))
	param.CreatedAt = dto.Sort(queryParams.Get("createdAt"))

	role := ctx.Value("userRole").(string)

	if role == "it" {
		h.service.GetUser(param).Send(ctx)
	} else {
		errs.NewUnauthorizedError("user is not authorized").Send(ctx)
		return
	}
}
