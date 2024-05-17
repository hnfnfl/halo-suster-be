package service

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"halo-suster/internal/pkg/middleware"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterUser(body model.User) errs.Response {
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

	// check NIP in database
	var count int
	if err = tx.QueryRow("SELECT COUNT(*) FROM users WHERE nip = $1", body.NIP).Scan(&count); err != nil {
		return errs.NewInternalError("query error", err)
	}

	if count > 0 {
		return errs.NewGenericError(http.StatusUnauthorized, "NIP already registered")
	}

	// insert user by role
	if body.Role == "it" {
		stmt := "INSERT INTO users (user_id, nip, password_hash, name, role) VALUES ($1, $2, $3, $4, $5)"
		if _, err = tx.Exec(stmt, body.UserID, body.NIP, body.PasswordHash, body.Name, body.Role); err != nil {
			return errs.NewInternalError("insert error", err)
		}
	} else if body.Role == "nurse" {
		stmt := "INSERT INTO users (user_id, nip, name, role, card_image) VALUES ($1, $2, $3, $4, $5)"
		if _, err = tx.Exec(stmt, body.UserID, body.NIP, body.Name, body.Role, body.CardImage); err != nil {
			return errs.NewInternalError("insert error", err)
		}
	}

	// generate token
	var token string
	if body.Role == "it" {
		token, err = middleware.JWTSign(s.Config(), body.UserID)
		if err != nil {
			return errs.NewInternalError("token signing error", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return errs.NewInternalError("commit error", err)
	}

	return errs.Response{
		Message: "User registered successfully",
		Data: dto.AuthResponse{
			UserId:      body.UserID,
			NIP:         body.NIP,
			Name:        body.Name,
			AccessToken: token,
		},
	}
}

func (s *Service) LoginUser(data model.User, body dto.LoginRequest) errs.Response {
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

	// check NIP in database
	q := "SELECT user_id, name, password_hash FROM users WHERE nip = $1"

	query_err := tx.QueryRow(q, data.NIP).Scan(&data.UserID, data.Name, data.PasswordHash)
	if query_err != nil {
		return errs.NewNotFoundError("user is not found")
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword([]byte(body.Password), []byte(data.PasswordHash)); err != nil {
		return errs.NewBadRequestError("password is wrong", err)
	}

	// generate token
	var token string
	token, err = middleware.JWTSign(s.Config(), data.UserID)
	if err != nil {
		return errs.NewInternalError("token signing error", err)
	}

	if err = tx.Commit(); err != nil {
		return errs.NewInternalError("commit error", err)
	}

	return errs.Response{
		Message: "User login successfully",
		Data: dto.AuthResponse{
			UserId:      data.UserID,
			NIP:         data.NIP,
			Name:        data.Name,
			AccessToken: token,
		},
	}
}
