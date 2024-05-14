package middleware

import (
	"halo-suster/internal/pkg/configuration"
	"halo-suster/internal/pkg/errs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JWTSign(cfg *configuration.Configuration, expiry time.Duration, userId string) (string, error) {
	timeStamp := time.Now()
	expiryTime := timeStamp.Add(expiry)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeStamp),
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			Subject:   userId,
		},
	)

	return token.SignedString([]byte(cfg.JWTSecret))
}

func JWTVerify(cfg *configuration.Configuration, tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrInvalidSigningMethod
		}

		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errs.ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", errs.ErrInvalidClaimsType
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return "", errs.ErrTokenExpired
	}

	return claims.Subject, nil
}