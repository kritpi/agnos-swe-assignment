package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kritpi/agnos-swe-assignment/handler/dto"
	"github.com/kritpi/agnos-swe-assignment/middleware"
)

// PatientSearch godoc
// @Summary      Search a patient
// @Description  Finds the one patient registered at the logged-in staff member's hospital that matches the filters. Every filter is optional, but at least one is required; filters are combined with AND.
// @Description  When nothing is stored locally and national_id or passport_id is given, the patient is fetched from the hospital's HIS, replaces the stored record, and is returned if it matches every filter.
// @Tags         patient
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.SearchPatientRequest  false  "Search filters"
// @Success      200   {object}  dto.PatientResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      409   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Failure      502   {object}  ErrorResponse
// @Failure      504   {object}  ErrorResponse
// @Router       /patient/search [post]
func (h *handler) PatientSearch(c *gin.Context) {
	ctx := c.Request.Context()

	claims, ok := middleware.StaffClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Code: http.StatusUnauthorized, Message: "missing staff claims"})
		return
	}

	// An empty body means no filters, which the service rejects.
	var req dto.SearchPatientRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	// The hospital always comes from the token, never from the body.
	criteria, err := req.ToDomain(claims.HospitalID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	patient, err := h.svc.SearchPatient(ctx, claims.HospitalCode, criteria)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.NewPatientResponse(*patient))
}
