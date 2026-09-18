package dto

import (
	"time"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

// CreateStaffRequest is the body of POST /staff/create.
type CreateStaffRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Hospital string `json:"hospital" binding:"required"`
}

// ToDomain converts the request to domain credentials.
func (r CreateStaffRequest) ToDomain() domain.StaffCredentials {
	return domain.StaffCredentials{
		Username:     r.Username,
		Password:     r.Password,
		HospitalCode: r.Hospital,
	}
}

// LoginRequest is the body of POST /staff/login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

// ToDomain converts the request to domain credentials.
func (r LoginRequest) ToDomain() domain.StaffCredentials {
	return domain.StaffCredentials{
		Username:     r.Username,
		Password:     r.Password,
		HospitalCode: r.Hospital,
	}
}

// StaffResponse is a staff member as returned by the API.
type StaffResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	HospitalID string `json:"hospital_id"`
	CreatedAt  string `json:"created_at"`
}

// NewStaffResponse builds a StaffResponse from the domain model.
func NewStaffResponse(s domain.Staff) StaffResponse {
	return StaffResponse{
		ID:         s.ID,
		Username:   s.Username,
		HospitalID: s.HospitalID,
		CreatedAt:  s.CreatedAt.Format(time.RFC3339),
	}
}

// LoginResponse is returned by a successful login.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
