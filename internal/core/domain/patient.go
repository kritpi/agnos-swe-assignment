package domain

import (
	"time"

	"github.com/kritpi/agnos-swe-assignment/enum"
)

// Patient is a person's identity, shared across hospitals. Empty strings mean
// the value is not set.
type Patient struct {
	ID           string
	NationalID   string
	PassportID   string
	FirstNameTH  string
	MiddleNameTH string
	LastNameTH   string
	FirstNameEN  string
	MiddleNameEN string
	LastNameEN   string
	DateOfBirth  time.Time
	Gender       enum.Gender
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// HospitalPatient is a patient as registered at a specific hospital.
type HospitalPatient struct {
	ID          string
	HospitalID  string
	PatientID   string
	PatientHN   string
	PhoneNumber string
	Email       string
	Patient     Patient
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PatientSearchCriteria filters patients within a hospital. Empty strings and
// a nil DateOfBirth mean the filter is not applied.
type PatientSearchCriteria struct {
	HospitalID  string
	NationalID  string
	PassportID  string
	FirstName   string
	MiddleName  string
	LastName    string
	DateOfBirth *time.Time
	PhoneNumber string
	Email       string
}
