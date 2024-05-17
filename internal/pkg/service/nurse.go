package service

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/errs"

	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) AccessNurse(Nurse model.User) errs.Response {
	var err error

	tx, err := s.DB().Begin()
	if err != nil {
		return errs.NewInternalError("transaction error", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}()

	stmt := "UPDATE users SET password = $1 where user_id = $2"
	_, queryErr := tx.Exec(stmt, Nurse.PasswordHash, Nurse.UserID)

	if queryErr != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return errs.NewNotFoundError("user id not found")
			}
		}
		return errs.NewInternalError("insert error", err)
	}

	return errs.Response{
		Code:    http.StatusOK,
		Message: "Nurse access successfully added",
	}
}
