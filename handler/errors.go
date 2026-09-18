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
	{domain.ErrHospitalNotFound, http.StatusNotFound},
	{domain.ErrNotFound, http.StatusNotFound},
	{domain.ErrStaffAlreadyExists, http.StatusConflict},
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
