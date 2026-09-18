package dto

import (
	"time"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

type CreateStaffRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Hospital string `json:"hospital" binding:"required"`
}

func (r CreateStaffRequest) ToDomain() domain.StaffCredentials {
	return domain.StaffCredentials{
		Username:     r.Username,
		Password:     r.Password,
		HospitalCode: r.Hospital,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

func (r LoginRequest) ToDomain() domain.StaffCredentials {
	return domain.StaffCredentials{
		Username:     r.Username,
		Password:     r.Password,
		HospitalCode: r.Hospital,
	}
}

type StaffResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	HospitalID string `json:"hospital_id"`
	CreatedAt  string `json:"created_at"`
}

func (s StaffResponse) FromDomain(dm *domain.Staff) *StaffResponse {
	return &StaffResponse{
		ID:         dm.ID,
		Username:   dm.Username,
		HospitalID: dm.HospitalID,
		CreatedAt:  dm.CreatedAt.Format(time.RFC3339),
	}
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (l LoginResponse) FromDomain(dm *domain.LoginResponse) *LoginResponse {
	return &LoginResponse{
		AccessToken: dm.AccessToken,
	}
}
