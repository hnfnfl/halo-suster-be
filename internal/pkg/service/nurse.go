package service

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"

	"net/http"
)

func (s *Service) AccessNurse(Nurse model.User) errs.Response {
	db := s.DB()

	data, errNotFound := s.FindUserById(Nurse.UserID, "nurse")
	if errNotFound.Error != "" {
		return errNotFound
	}

	if data.Role != "nurse" {
		return errs.NewNotFoundError("user is not a nurse (nip not starts with 303)", errs.ErrUserNotFound)
	}

	stmt := "UPDATE users SET password_hash = $1 where user_id = $2"
	_, queryErr := db.Exec(stmt, Nurse.PasswordHash, Nurse.UserID)

	if queryErr != nil {
		return errs.NewInternalError("update password error", queryErr)
	}

	return errs.Response{
		Code:    http.StatusOK,
		Message: "Nurse access successfully added",
	}
}

func (s *Service) UpdateNurse(req dto.RequestUpdateNurse, userId string) errs.Response {
	db := s.DB()

	data, errNotFound := s.FindUserById(userId, "nurse")
	if errNotFound.Error != "" {
		return errNotFound
	}

	if data.Role != "nurse" {
		return errs.NewNotFoundError("user is not a nurse (nip not starts with 303)", errs.ErrUserNotFound)
	}

	err := s.FindExistingNIP(req.NIP, "nurse")
	if err != nil {
		return errs.NewGenericError(http.StatusConflict, "nip is exists")
	}

	stmt := "UPDATE users SET nip = $1, name = $2 where user_id = $3"
	_, queryErr := db.Exec(stmt, req.NIP, req.Name, userId)

	if queryErr != nil {
		return errs.NewInternalError("update users error", queryErr)
	}

	return errs.Response{
		Code:    http.StatusOK,
		Message: "Nurse data successfully updated",
	}
}

func (s *Service) DeleteNurse(userId string) errs.Response {
	db := s.DB()

	data, errNotFound := s.FindUserById(userId, "nurse")
	if errNotFound.Error != "" {
		return errNotFound
	}

	if data.Role != "nurse" {
		return errs.NewNotFoundError("user is not a nurse (nip not starts with 303)", errs.ErrUserNotFound)
	}

	stmt := "DELETE FROM users where role = 'nurse' and user_id = $1"
	_, queryErr := db.Exec(stmt, userId)

	if queryErr != nil {
		return errs.NewInternalError("delete users error", queryErr)
	}

	return errs.Response{
		Code:    http.StatusOK,
		Message: "Nurse data successfully updated",
	}
}
