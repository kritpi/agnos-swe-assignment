package port

import (
	"context"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

type Service interface {
	CreateStaff(ctx context.Context, cred domain.StaffCredentials) (*domain.Staff, error)
	LoginStaff(ctx context.Context, cred domain.StaffCredentials) (*domain.LoginResponse, error)
	VerifyStaff(ctx context.Context, hospitalID, staffID string) error
	SearchPatient(ctx context.Context, hospitalCode string, c domain.PatientSearchCriteria) (*domain.HospitalPatient, error)
}
