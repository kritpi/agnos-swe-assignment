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
	FindStaffByID(ctx context.Context, hospitalID, staffID string) (*domain.Staff, error)
	CreateStaff(ctx context.Context, s domain.Staff) error

	// Patient
	SearchPatients(ctx context.Context, c domain.PatientSearchCriteria) ([]domain.HospitalPatient, error)
	FindPatientByNationalID(ctx context.Context, nationalID string) (*domain.Patient, error)
	FindPatientByPassportID(ctx context.Context, passportID string) (*domain.Patient, error)
	CreatePatient(ctx context.Context, p domain.Patient) error
	UpdatePatient(ctx context.Context, p domain.Patient) error
	UpsertHospitalPatient(ctx context.Context, hp domain.HospitalPatient) (*domain.HospitalPatient, error)
}
