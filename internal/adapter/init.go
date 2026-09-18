// Package adapter talks to the hospital information systems (HIS). Each
// hospital has its own operation file, because every HIS has its own URL,
// request shape and payload.
package adapter

import (
	"net/http"

	"github.com/kritpi/agnos-swe-assignment/internal/core/port"
	"github.com/kritpi/agnos-swe-assignment/property"
)

// adapter is the concrete HIS client. The struct and its constructor live
// here; each operation is a receiver method in its own file.
type adapter struct {
	cfg *property.Config
	// httpClient has no timeout of its own: each operation applies its
	// hospital's timeout through the request context.
	httpClient *http.Client
}

// New builds an Adapter from the application configuration.
func New(cfg *property.Config) port.Adapter {
	return &adapter{
		cfg:        cfg,
		httpClient: &http.Client{},
	}
}
