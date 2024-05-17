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

	data, errNotFound := s.FindUserById(Nurse.UserID)
	if errNotFound.Error != nil {
		return errNotFound
	}

	if data.Role != "nurse" {
		return errs.NewNotFoundError("user is not a nurse (nip not starts with 303)", errs.ErrUnauthorized)
	}

	stmt := "UPDATE users SET password = $1 where user_id = $2"
	_, queryErr := tx.Exec(stmt, Nurse.PasswordHash, Nurse.UserID)

	if queryErr != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return errs.NewNotFoundError("user id not found", pgErr)
			}
		}
		return errs.NewInternalError("insert error", err)
	}

	if queryErr = tx.Commit(); queryErr != nil {
		return errs.NewInternalError("commit error", queryErr)
	}

	return errs.Response{
		Code:    http.StatusOK,
		Message: "Nurse access successfully added",
	}
}
