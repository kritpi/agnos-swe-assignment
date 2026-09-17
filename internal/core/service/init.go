package service

import (
	"github.com/kritpi/agnos-swe-assignment/internal/core/port"
)

// service is the concrete business-logic implementation. The struct and its
// constructor live here; each operation is a receiver method in its own file.
type service struct {
	repo port.Repository
}

// New builds a Service from a Repository.
func New(repo port.Repository) port.Service {
	return &service{
		repo: repo,
	}
}
