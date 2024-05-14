package errs

import "errors"

// nolint::lll
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPasswordMissmatch = errors.New("password missmatch")
)
