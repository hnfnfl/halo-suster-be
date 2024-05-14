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

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	// set db to gin context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// t, err := middleware.JWTSign(cfg, 10*time.Minute, "1234567890")
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }
	// log.Infof("Token: %v", t)

	// test, err := middleware.JWTVerify(cfg, t)
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }
	// log.Infof("Test: %v", test)

	// login
	// loginGroup := router.Group("/v1/user/")
	// loginGroup.POST("/it/login", Login)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router.Run(":8080")
}
