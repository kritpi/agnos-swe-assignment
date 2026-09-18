package service

import (
	"github.com/kritpi/agnos-swe-assignment/internal/core/port"
	"github.com/kritpi/agnos-swe-assignment/property"
)

type service struct {
	repo     port.Repository
	adapters port.Adapter
	cfg      *property.Config
}

func New(repo port.Repository, adapters port.Adapter, cfg *property.Config) port.Service {
	return &service{
		repo:     repo,
		adapters: adapters,
		cfg:      cfg,
	}
}
