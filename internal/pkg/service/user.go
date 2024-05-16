package service

import (
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"halo-suster/internal/pkg/middleware"
	"net/http"
	"time"
)

func (s *Service) RegisterUser(in model.User) errs.Response {
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
	if err = tx.QueryRow("SELECT COUNT(*) FROM users WHERE nip = $1", in.NIP).Scan(&count); err != nil {
		return errs.NewInternalError("query error", err)
	}

	if count > 0 {
		return errs.NewGenericError(http.StatusUnauthorized, "NIP already registered")
	}

	// insert user by role
	if in.Role == "it" {
		stmt := "INSERT INTO users (user_id, nip, password_hash, name, role) VALUES ($1, $2, $3, $4, $5)"
		if _, err = tx.Exec(stmt, in.UserID, in.NIP, in.PasswordHash, in.Name, in.Role); err != nil {
			return errs.NewInternalError("insert error", err)
		}
	} else if in.Role == "nurse" {
		stmt := "INSERT INTO users (user_id, nip, name, role, card_image) VALUES ($1, $2, $3, $4, $5)"
		if _, err = tx.Exec(stmt, in.UserID, in.NIP, in.Name, in.Role, in.CardImage); err != nil {
			return errs.NewInternalError("insert error", err)
		}
	}

	// generate token
	var token string
	if in.Role == "it" {
		token, err = middleware.JWTSign(s.Config(), 1*time.Hour, in.UserID)
		if err != nil {
			return errs.NewInternalError("token signing error", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return errs.NewInternalError("commit error", err)
	}

	return errs.Response{
		Message: "User registered successfully",
		Data: dto.RegisterOutput{
			UserId:      in.UserID,
			NIP:         in.NIP,
			Name:        in.Name,
			AccessToken: token,
		},
	}
}
