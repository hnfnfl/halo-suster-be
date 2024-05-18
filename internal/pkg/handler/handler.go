package handler

import (
	"fmt"
	"halo-suster/internal/db"
	"halo-suster/internal/pkg/configuration"
	"halo-suster/internal/pkg/middleware"
	"halo-suster/internal/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	service := service.NewService(cfg, db)
	userHandler := NewUserHandler(service, &validator.Validate{})
	patientHandler := NewPatientHandler(service)
	nurseHandler := NewNurseHandler(service, &validator.Validate{})

	// login
	authGroup := router.Group("/v1/user/")
	authGroup.POST("it/register", userHandler.Register)
	authGroup.POST("it/login", userHandler.Login)
	authGroup.POST("nurse/login", userHandler.Login)

	nurseGroup := router.Group("/v1/user/")
	nurseGroup.Use(middleware.JWTAuthMiddleware(cfg))
	nurseGroup.POST("nurse/register", userHandler.Register)
	nurseGroup.POST("nurse/{userId}/access", nurseHandler.AccessNurse)
	// nurseGroup.GET("nurse/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	medicalRecord := router.Group("/v1/medical/")
	medicalRecord.Use(middleware.JWTAuthMiddleware(cfg))
	medicalRecord.POST("patient", patientHandler.CreatePatient)
	medicalRecord.GET("patient", patientHandler.GetPatient)

	return router.Run(":8080")
}
