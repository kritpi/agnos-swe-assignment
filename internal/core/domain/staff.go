package domain

import "time"

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
