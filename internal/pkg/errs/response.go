package errs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	respCode int

	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewGenericError(code int, msg string) Response {
	return Response{
		respCode: code,
		Message:  msg,
	}
}

func NewInternalError(msg string, err error) Response {
	return Response{
		respCode: http.StatusInternalServerError,
		Error:    err.Error(),
		Message:  msg,
	}
}

func NewValidationError(msg string, err error) Response {
	return Response{
		respCode: http.StatusBadRequest,
		Error:    err.Error(),
		Message:  msg,
	}
}

func NewUnauthorizedError(msg string) Response {
	return Response{
		respCode: http.StatusUnauthorized,
		Message:  msg,
	}
}

func (e Response) Send(ctx *gin.Context) {
	ctx.JSON(e.respCode, e)
}
