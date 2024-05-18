package service

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"

	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) AccessNurse(Nurse model.User) errs.Response {
	var err error

	db := s.DB()

	data, errNotFound := s.FindUserById(Nurse.UserID)
	if errNotFound.Error != "" {
		return errNotFound
	}

	if data.Role != "nurse" {
		return errs.NewNotFoundError("user is not a nurse (nip not starts with 303)", errs.ErrUserNotFound)
	}

	stmt := "UPDATE users SET password_hash = $1 where user_id = $2"
	_, queryErr := db.Exec(stmt, Nurse.PasswordHash, Nurse.UserID)

	if queryErr != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return errs.NewNotFoundError("user id not found", pgErr)
			}
		}
		return errs.NewInternalError("update password error", err)
	}

	return errs.Response{
		Code:    http.StatusOK,
		Message: "Nurse access successfully added",
	}
}

func (s *Service) UpdateNurse(req dto.RequestUpdateNurse, userId string) errs.Response {
	var err error

	db := s.DB()

	data, errNotFound := s.FindUserById(userId)
	if errNotFound.Error != "" {
		return errNotFound
	}

	if data.Role != "nurse" {
		return errs.NewNotFoundError("user is not a nurse (nip not starts with 303)", errs.ErrUserNotFound)
	}

	stmt := "UPDATE users SET nip = $1, name = $2 where user_id = $2"
	_, queryErr := db.Exec(stmt, req.NIP, req.Name, userId)

	if queryErr != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errs.NewGenericError(http.StatusConflict, "nip is exists")
			} else if pgErr.Code == "22P02" {
				return errs.NewNotFoundError("user id not found", pgErr)
			}
		}
		return errs.NewInternalError("update users error", err)
	}

	return errs.Response{
		Code:    http.StatusOK,
		Message: "Nurse data successfully updated",
	}
}

func (s *Service) DeleteNurse(userId string) errs.Response {
	var err error

	db := s.DB()

	data, errNotFound := s.FindUserById(userId)
	if errNotFound.Error != "" {
		return errNotFound
	}

	if data.Role != "nurse" {
		return errs.NewNotFoundError("user is not a nurse (nip not starts with 303)", errs.ErrUserNotFound)
	}

	stmt := "DELET FROM users where user_id = $1"
	_, queryErr := db.Exec(stmt, userId)

	if queryErr != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return errs.NewNotFoundError("user id not found", pgErr)
			}
		}
		return errs.NewInternalError("update users error", err)
	}

	return errs.Response{
		Code:    http.StatusOK,
		Message: "Nurse data successfully updated",
	}
}
