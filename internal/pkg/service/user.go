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

func (s *Service) LoginUser(body model.User) errs.Response {
	var (
		err error
		out model.User
	)

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
	stmt := "SELECT user_id, nip, name, password_hash FROM users WHERE nip = $1"
	if err := tx.QueryRow(stmt, body.NIP).Scan(
		&out.UserID,
		&out.NIP,
		&out.Name,
		&out.PasswordHash,
	); err != nil {
		return errs.NewNotFoundError("user is not found")
	}

	if out.PasswordHash == nil {
		return errs.NewBadRequestError("user is not having access", errs.ErrUnauthorized)
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword(out.PasswordHash, body.PasswordHash); err != nil {
		return errs.NewBadRequestError("password is wrong", err)
	}

	// generate token
	var token string
	token, err = middleware.JWTSign(s.Config(), body.UserID)
	if err != nil {
		return errs.NewInternalError("token signing error", err)
	}

	if err = tx.Commit(); err != nil {
		return errs.NewInternalError("commit error", err)
	}

	return errs.Response{
		Message: "User login successfully",
		Data: dto.AuthResponse{
			UserId:      out.UserID,
			NIP:         out.NIP,
			Name:        out.Name,
			AccessToken: token,
		},
	}
}
