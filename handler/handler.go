package handler

import (
	"github.com/kritpi/agnos-swe-assignment/internal/core/port"
)

// Handler is the HTTP handler interface. The struct and constructor live here;
// each endpoint is implemented as a receiver method in its own file.
type Handler interface {
}

type handler struct {
	svc port.Service
}

// New builds a Handler from a Service.
func New(svc port.Service) Handler {
	return &handler{svc: svc}
}

// ErrorResponse is a generic error envelope, referenced by Swagger @Failure annotations.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
