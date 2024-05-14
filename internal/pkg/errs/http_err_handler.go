package errs

type ErrorRespone struct {
	ErrorType string `json:"error_type"`
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_message"`
}

func NewInternalError() ErrorRespone {
	return ErrorRespone{
		ErrorType: "InternalError",
		ErrorCode: 500,
		ErrorMsg:  "Internal server error",
	}
}

func NewValidationError(code int, msg string) ErrorRespone {
	return ErrorRespone{
		ErrorType: "ValidationError",
		ErrorCode: code,
		ErrorMsg:  msg,
	}
}