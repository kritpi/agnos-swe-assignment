package port

import (
	"context"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

type Repository interface {
	Transactional(ctx context.Context, f func(ctx context.Context) error) error

	// Hospital
	FindHospitalByCode(ctx context.Context, code string) (*domain.Hospital, error)

	// Staff
	FindStaffByUsername(ctx context.Context, hospitalID, username string) (*domain.Staff, error)
	CreateStaff(ctx context.Context, s domain.Staff) error
}
