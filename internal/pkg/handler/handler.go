package handler

import (
	"fmt"
	"halo-suster/internal/db"
	"halo-suster/internal/pkg/configuration"
	"halo-suster/internal/pkg/middleware"
	"halo-suster/internal/pkg/service"

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

	service := service.NewService(cfg, db)
	userHandler := NewUserHandler(service)
	patientHandler := NewPatientHandler(service)
	nurseHandler := NewNurseHandler(service)
	medicalRecordHandler := NewMedicalRecordHandler(service)
	imageHandler := NewImageHandler(service)

	// login
	authGroup := router.Group("/v1/user/")
	authGroup.POST("it/register", userHandler.Register)
	authGroup.POST("it/login", userHandler.Login)
	authGroup.POST("nurse/login", userHandler.Login)

	nurseGroup := router.Group("/v1/user/")
	nurseGroup.Use(middleware.JWTAuthMiddleware(cfg))
	nurseGroup.POST("nurse/register", userHandler.Register)
	nurseGroup.GET("", userHandler.GetUsers)
	nurseGroup.PUT("nurse/:userId", nurseHandler.UpdateNurse)
	nurseGroup.DELETE("nurse/:userId", nurseHandler.DeleteNurse)
	nurseGroup.POST("nurse/:userId/access", nurseHandler.AccessNurse)

	imageUpload := router.Group("/v1/image/")
	imageUpload.Use(middleware.JWTAuthMiddleware(cfg))
	imageUpload.POST("/", imageHandler.UploadImage)
	// nurseGroup.GET("nurse/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	medicalRecord := router.Group("/v1/medical/")
	medicalRecord.Use(middleware.JWTAuthMiddleware(cfg))
	medicalRecord.POST("patient", patientHandler.CreatePatient)
	medicalRecord.GET("patient", patientHandler.GetPatient)
	medicalRecord.POST("record", medicalRecordHandler.CreateMedicalRecord)
	medicalRecord.GET("record", medicalRecordHandler.GetMedicalRecord)

	return router.Run(":8080")
}
