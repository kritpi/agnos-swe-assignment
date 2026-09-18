package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kritpi/agnos-swe-assignment/internal/core/port"
)

type Handler interface {
	StaffCreate(c *gin.Context)
	StaffLogin(c *gin.Context)
	PatientSearch(c *gin.Context)
}

type handler struct {
	svc port.Service
}

func New(svc port.Service) Handler {
	return &handler{svc: svc}
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
