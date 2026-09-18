package service

import (
	"github.com/kritpi/agnos-swe-assignment/internal/core/port"
)

// service is the concrete business-logic implementation. The struct and its
// constructor live here; each operation is a receiver method in its own file.
type service struct {
	repo     port.Repository
	adapters port.Adapter
}

// New builds a Service from a Repository and the HIS adapters, keyed by
// hospital code.
func New(repo port.Repository, adapters port.Adapter) port.Service {
	return &service{
		repo:     repo,
		adapters: adapters,
	}
}
