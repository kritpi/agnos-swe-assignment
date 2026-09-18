package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kritpi/agnos-swe-assignment/handler/dto"
)

// StaffCreate godoc
// @Summary      Create a staff member
// @Description  Creates a staff account for the given hospital.
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateStaffRequest  true  "Staff credentials"
// @Success      201   {object}  dto.StaffResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      409   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /staff/create [post]
func (h *handler) StaffCreate(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	staff, err := h.svc.CreateStaff(ctx, req.ToDomain())
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.StaffResponse{}.FromDomain(staff))
}
