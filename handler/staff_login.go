package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kritpi/agnos-swe-assignment/handler/dto"
)

// StaffLogin godoc
// @Summary      Log in as a staff member
// @Description  Checks the staff credentials for the given hospital and returns a JWT access token.
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest  true  "Staff credentials"
// @Success      200   {object}  dto.LoginResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /staff/login [post]
func (h *handler) StaffLogin(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	accessToken, err := h.svc.LoginStaff(ctx, req.ToDomain())
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{}.FromDomain(accessToken))
}
