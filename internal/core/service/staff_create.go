package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

func (s *service) CreateStaff(ctx context.Context, in domain.StaffCredentials) (*domain.Staff, error) {
	// check hospital exists
	hospital, err := s.repo.FindHospitalByCode(ctx, in.HospitalCode)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrHospitalNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find hospital: %w", err)
	}

	// check username exists at this hospital (UNIQUE (hospital_id, username))
	_, err = s.repo.FindStaffByUsername(ctx, hospital.ID, in.Username)
	if err == nil {
		return nil, domain.ErrStaffAlreadyExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("find staff: %w", err)
	}

	// hash password 
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// staff uuid
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate staff id: %w", err)
	}

	now := time.Now().UTC()
	staff := domain.Staff{
		ID:           id.String(),
		HospitalID:   hospital.ID,
		Username:     in.Username,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// create new staff
	if err := s.repo.CreateStaff(ctx, staff); err != nil {
		if errors.Is(err, domain.ErrStaffAlreadyExists) {
			return nil, err
		}
		return nil, fmt.Errorf("create staff: %w", err)
	}

	return &staff, nil
}
