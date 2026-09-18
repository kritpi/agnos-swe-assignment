package entity

import (
	"time"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

type HospitalPatient struct {
	ID          string    `db:"id"`
	HospitalID  string    `db:"hospital_id"`
	PatientID   string    `db:"patient_id"`
	PatientHN   string    `db:"patient_hn"`
	PhoneNumber string    `db:"phone_number"`
	Email       string    `db:"email"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (e HospitalPatient) ToDomain() domain.HospitalPatient {
	return domain.HospitalPatient{
		ID:          e.ID,
		HospitalID:  e.HospitalID,
		PatientID:   e.PatientID,
		PatientHN:   e.PatientHN,
		PhoneNumber: e.PhoneNumber,
		Email:       e.Email,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func HospitalPatientFromDomain(d domain.HospitalPatient) HospitalPatient {
	return HospitalPatient{
		ID:          d.ID,
		HospitalID:  d.HospitalID,
		PatientID:   d.PatientID,
		PatientHN:   d.PatientHN,
		PhoneNumber: d.PhoneNumber,
		Email:       d.Email,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
