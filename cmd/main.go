package main

import (
	"halo-suster/internal/pkg/configuration"
	"halo-suster/internal/pkg/handler"
	"halo-suster/internal/pkg/logger"
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

	if err := handler.Run(config, log); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
