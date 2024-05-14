package router

import (
	"fmt"
	"halo-suster/internal/db"
	"halo-suster/internal/pkg/configuration"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run(cfg *configuration.Configuration, log *logrus.Logger) error {
	db, err := db.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	router := gin.Default()
	// set db to gin context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// login
	// loginGroup := router.Group("/v1/user/")
	// loginGroup.POST("/it/login", Login)

	return nil
}
