package main

import (
	"halo-suster/internal/pkg/configuration"
	"halo-suster/internal/pkg/logger"
	"halo-suster/internal/pkg/router"
)

func main() {
	config, err := configuration.NewConfiguration()
	if err != nil {
		panic(err)
	}

	log, err := logger.NewLogger(config.LogLevel)
	if err != nil {
		panic(err)
	}

	if err := router.Run(config, log); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
