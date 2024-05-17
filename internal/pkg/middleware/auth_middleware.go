package middleware

import (
	"halo-suster/internal/pkg/configuration"
	"halo-suster/internal/pkg/errs"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func JWTSign(cfg *configuration.Configuration, userId string, nip string, role string) (string, error) {
	expiry := time.Duration(cfg.AuthExpiry) * time.Hour
	timeStamp := time.Now()
	expiryTime := timeStamp.Add(expiry)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeStamp),
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			ID:        userId,
			Issuer:    nip,
			Subject:   role,
		},
	)

	return token.SignedString([]byte(cfg.JWTSecret))
}

func JWTAuthMiddleware(cfg *configuration.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := ""

		if authHeader != "" {
			tokenString = authHeader[len(BearerSchema):]
		}

		if tokenString == "" {
			errs.NewUnauthorizedError("Authorization header not provided").Send(c)
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errs.ErrInvalidSigningMethod
			}

			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			errs.NewUnauthorizedError(err.Error()).Send(c)
			c.Abort()
			return
		}

		if !token.Valid {
			errs.NewUnauthorizedError(errs.ErrInvalidToken.Error()).Send(c)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			errs.NewUnauthorizedError(errs.ErrInvalidClaimsType.Error()).Send(c)
			c.Abort()
			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			errs.NewUnauthorizedError(errs.ErrTokenExpired.Error()).Send(c)
			c.Abort()
			return
		}

		c.Set("userID", claims.Subject)

		c.Next()
	}
}

func PasswordHash(password string, salt string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
}
