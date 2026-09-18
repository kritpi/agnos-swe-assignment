package domain

import (
	"strings"
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

// IsEmpty reports whether no filter is set. HospitalID is not a filter: it
// always comes from the staff member's token.
func (c PatientSearchCriteria) IsEmpty() bool {
	return !c.HasPatientID() && c.FirstName == "" && c.MiddleName == "" && c.LastName == "" &&
		c.DateOfBirth == nil && c.PhoneNumber == "" && c.Email == ""
}

func (c PatientSearchCriteria) HasPatientID() bool {
	return c.NationalID != "" || c.PassportID != ""
}

func (c PatientSearchCriteria) Matches(hp HospitalPatient) bool {
	p := hp.Patient
	switch {
	case c.NationalID != "" && c.NationalID != p.NationalID,
		c.PassportID != "" && c.PassportID != p.PassportID,
		!matchesName(c.FirstName, p.FirstNameTH, p.FirstNameEN),
		!matchesName(c.MiddleName, p.MiddleNameTH, p.MiddleNameEN),
		!matchesName(c.LastName, p.LastNameTH, p.LastNameEN),
		c.DateOfBirth != nil && !sameDate(*c.DateOfBirth, p.DateOfBirth),
		c.PhoneNumber != "" && c.PhoneNumber != hp.PhoneNumber,
		c.Email != "" && !strings.EqualFold(c.Email, hp.Email):
		return false
	}
	return true
}

func matchesName(filter, th, en string) bool {
	return filter == "" || strings.EqualFold(filter, th) || strings.EqualFold(filter, en)
}

func sameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
