package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Staff struct {
	ID           string
	HospitalID   string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Credential for login
type StaffCredentials struct {
	Username     string
	Password     string
	HospitalCode string
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type StaffClaims struct {
	StaffID      string `json:"staff_id"`
	Username     string `json:"username"`
	HospitalID   string `json:"hospital_id"`
	HospitalCode string `json:"hospital_code"`
	jwt.RegisteredClaims
}
