package service

import (
	"database/sql"
	"fmt"
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

func (s *Service) UpdateNurse(userId string, body dto.RequestUpdateNurse) errs.Response {
	var changedRow model.User
	db := s.DB()

	// count the number of nurses with the same NIP
	q := "SELECT COUNT(*) FROM users WHERE nip = $1 AND role = 'nurse'"
	var count int
	if err := db.QueryRow(q, body.NIP).Scan(&count); err != nil {
		return errs.Response{
			Code:    http.StatusInternalServerError,
			Message: "query error",
		}
	}

	if count > 0 {
		return errs.Response{
			Code:    http.StatusConflict,
			Message: "NIP already exists",
		}
	}

	q = "UPDATE users SET nip = $1, name = $2 WHERE user_id = $3 AND role = 'nurse' RETURNING user_id"
	if err := db.QueryRow(q, body.NIP, body.Name, userId).Scan(&changedRow.UserID); err != nil {
		if err == sql.ErrNoRows {
			return errs.Response{
				Code:    http.StatusNotFound,
				Message: "User not found",
			}
		} else {
			return errs.Response{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("Failed to update nurse: %v", err),
			}
		}
	}

	// If no errors, return a success response
	return errs.Response{
		Code:    http.StatusOK,
		Message: "Nurse updated successfully",
	}
}

func (s *Service) DeleteNurse(userId string) errs.Response {
	db := s.DB()
	q := "DELETE FROM users WHERE user_id = $1 AND role = 'nurse' RETURNING user_id"

	var deletedRow model.User

	queryErr := db.QueryRow(q, userId).Scan(&deletedRow.UserID)
	if queryErr == sql.ErrNoRows {
		return errs.Response{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}

	if queryErr != nil {
		// Handle other types of errors
		return errs.Response{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Failed to delete nurse: %v", queryErr),
		}
	}

	// If no errors, return a success response
	return errs.Response{
		Code:    http.StatusOK,
		Message: "Nurse deleted successfully",
	}
}
