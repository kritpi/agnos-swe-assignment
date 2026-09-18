package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

var errorStatus = []struct {
	err    error
	status int
}{
	{domain.ErrEmptySearchCriteria, http.StatusBadRequest},
	{domain.ErrHospitalNotFound, http.StatusNotFound},
	{domain.ErrPatientNotFound, http.StatusNotFound},
	{domain.ErrNotFound, http.StatusNotFound},
	{domain.ErrStaffAlreadyExists, http.StatusConflict},
	{domain.ErrInvalidCredentials, http.StatusUnauthorized},
	{domain.ErrPatientAmbiguous, http.StatusConflict},
	{domain.ErrPatientConflict, http.StatusConflict},
	{domain.ErrPatientAlreadyExists, http.StatusConflict},
	{domain.ErrHISTimeout, http.StatusGatewayTimeout},
	{domain.ErrHISUnavailable, http.StatusBadGateway},
	{domain.ErrHISInvalidResponse, http.StatusBadGateway},
}

func writeError(c *gin.Context, err error) {
	for _, m := range errorStatus {
		if errors.Is(err, m.err) {
			c.JSON(m.status, ErrorResponse{Code: m.status, Message: m.err.Error()})
			return
		}
	}
	_ = c.Error(err)
	c.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: "internal server error"})
}
