package service

import (
	"fmt"
	"halo-suster/internal/db/model"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"halo-suster/internal/pkg/middleware"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterUser(body model.User) errs.Response {
	var err error

	db := s.DB()

	// check NIP in database
	var count int
	if err = db.QueryRow("SELECT COUNT(*) FROM users WHERE nip = $1", body.NIP).Scan(&count); err != nil {
		return errs.NewInternalError("query error", err)
	}

	if count > 0 {
		return errs.NewGenericError(http.StatusUnauthorized, "NIP already registered")
	}

	// insert user by role
	if body.Role == "it" {
		stmt := "INSERT INTO users (user_id, nip, password_hash, name, role) VALUES ($1, $2, $3, $4, $5)"
		if _, err = db.Exec(stmt, body.UserID, body.NIP, body.PasswordHash, body.Name, body.Role); err != nil {
			return errs.NewInternalError("insert error", err)
		}
	} else if body.Role == "nurse" {
		stmt := "INSERT INTO users (user_id, nip, name, role, card_image) VALUES ($1, $2, $3, $4, $5)"
		if _, err = db.Exec(stmt, body.UserID, body.NIP, body.Name, body.Role, body.CardImage); err != nil {
			return errs.NewInternalError("insert error", err)
		}
	}

	// generate token
	var token string
	if body.Role == "it" {
		token, err = middleware.JWTSign(s.Config(), body.UserID, body.NIP, body.Role)
		if err != nil {
			return errs.NewInternalError("token signing error", err)
		}
	}

	return errs.Response{
		Code:    http.StatusCreated,
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

	db := s.DB()

	// check NIP in database
	stmt := "SELECT user_id, nip, name, password_hash, role FROM users WHERE nip = $1"
	if err := db.QueryRow(stmt, body.NIP).Scan(
		&out.UserID,
		&out.NIP,
		&out.Name,
		&out.PasswordHash,
		&out.Role,
	); err != nil {
		return errs.NewNotFoundError("user is not found", errs.ErrUserNotFound)
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
	token, err = middleware.JWTSign(s.Config(), out.UserID, out.NIP, out.Role)
	if err != nil {
		return errs.NewInternalError("token signing error", err)
	}

	return errs.Response{
		Code:    http.StatusOK,
		Message: "User login successfully",
		Data: dto.AuthResponse{
			UserId:      out.UserID,
			NIP:         out.NIP,
			Name:        out.Name,
			AccessToken: token,
		},
	}
}

func (s *Service) FindUserById(userId string) (model.User, errs.Response) {
	var err error

	db := s.DB()

	data := model.User{}

	// check NIP in database
	q := "SELECT user_id, nip, name, role FROM users WHERE user_id = $1"

	queryErr := db.QueryRow(q, userId).Scan(&data.UserID, &data.NIP, &data.Name, &data.Role)
	if queryErr != nil {
		return model.User{}, errs.NewInternalError("user is not found", err)
	}
	return data, errs.Response{}
}

func (s *Service) GetUser(param dto.ReqParamUserGet) errs.Response {
	var err error

	db := s.DB()

	var query strings.Builder

	query.WriteString(`SELECT user_id, nip, name, created_at FROM users WHERE 1=1 `)

	if param.UserID != "" {
		query.WriteString(fmt.Sprintf("AND user_id = '%s' ", param.UserID))
	}

	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	if param.NIP != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(nip) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	if param.Role != "" {
		switch dto.Role(param.Role) {
		case dto.IT:
			query.WriteString(fmt.Sprintf("AND role = '%s' ", param.Role))
		case dto.Nurse:
			query.WriteString(fmt.Sprintf("AND role = '%s' ", param.Role))
		default:
		}
	}

	if param.CreatedAt == "asc" {
		query.WriteString("ORDER BY created_at ASC ")
	} else {
		query.WriteString("ORDER BY created_at DESC ")
	}
	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	rows, err := db.Query(query.String())
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "02000" {
				return errs.Response{
					Code:    http.StatusOK,
					Message: "Get data successfully, but no data",
					Data:    []dto.ResUserGet{},
				}
			}
		}
		return errs.NewInternalError(err.Error(), err)
	}
	defer rows.Close()

	results := []dto.ResUserGet{}
	for rows.Next() {
		var createdAt time.Time
		result := dto.ResUserGet{}
		err := rows.Scan(
			&result.UserID,
			&result.NIP,
			&result.Name,
			&createdAt)
		if err != nil {
			return errs.NewInternalError(err.Error(), err)
		}
		result.CreatedAt = createdAt.Format(time.RFC3339)
		results = append(results, result)
	}
	return errs.Response{
		Code:    http.StatusOK,
		Message: "Get data successfully",
		Data:    results,
	}
}
