package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

func (s *service) LoginStaff(ctx context.Context, cred domain.StaffCredentials) (*domain.LoginResponse, error) {
	// check hospital exists
	hospital, err := s.repo.FindHospitalByCode(ctx, cred.HospitalCode)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("find hospital: %w", err)
	}

	// check staff exists at this hospital
	staff, err := s.repo.FindStaffByUsername(ctx, hospital.ID, cred.Username)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("find staff: %w", err)
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(cred.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// issue JWT
	now := time.Now().UTC()
	claims := domain.StaffClaims{
		StaffID:      staff.ID,
		Username:     staff.Username,
		HospitalID:   hospital.ID,
		HospitalCode: hospital.Code,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   staff.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.JWT.TTL)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	return &domain.LoginResponse{AccessToken: token}, nil
}
