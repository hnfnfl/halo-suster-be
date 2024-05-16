package service

import (
	"database/sql"
	"halo-suster/internal/pkg/configuration"
)

type Service struct {
	cfg *configuration.Configuration
	db  *sql.DB
}

func NewService(cfg *configuration.Configuration, db *sql.DB) *Service {
	return &Service{
		cfg: cfg,
		db:  db,
	}
}

func (s *Service) DB() *sql.DB {
	return s.db
}

func (s *Service) Config() *configuration.Configuration {
	return s.cfg
}
