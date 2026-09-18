package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrHospitalNotFound   = errors.New("hospital not found")
	ErrStaffAlreadyExists = errors.New("username already exists at this hospital")
	ErrInvalidCredentials = errors.New("invalid username, password or hospital")
)
