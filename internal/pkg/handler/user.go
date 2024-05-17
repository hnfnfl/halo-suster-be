package handler

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"halo-suster/internal/pkg/middleware"
	"halo-suster/internal/pkg/service"
	"halo-suster/internal/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	*service.Service
}

func NewUserHandler(s *service.Service) *UserHandler {
	return &UserHandler{s}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	in := dto.RegisterInput{}
	msg, err := util.JsonBinding(ctx, &in)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	// validate input
	if err := in.Validate(); err != nil {
		errs.NewValidationError("input validation error", err).Send(ctx)
		return
	}

	data := model.User{
		NIP:  strconv.Itoa(in.NIP),
		Name: in.Name,
	}

	var passHash []byte
	if in.Password != "" {
		var err error
		passHash, err = middleware.PasswordHash(in.Password, h.Config().Salt)
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

	h.RegisterUser(data).Send(ctx)
}
