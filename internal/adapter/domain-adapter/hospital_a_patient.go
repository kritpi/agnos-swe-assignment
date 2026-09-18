// Package domainadapter holds the hospitals' HIS wire formats. The field names
// and types are each hospital's, not ours; ToDomain converts them to the core
// domain model.
package domainadapter

import (
	"time"

	"github.com/kritpi/agnos-swe-assignment/enum"
	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

// HospitalAPatient is the response body of hospital A's patient search endpoint.
type HospitalAPatient struct {
	FirstNameTH  string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH   string `json:"last_name_th"`
	FirstNameEN  string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN   string `json:"last_name_en"`
	DateOfBirth  string `json:"date_of_birth"` // YYYY-MM-DD
	PatientHN    string `json:"patient_hn"`
	NationalID   string `json:"national_id"`
	PassportID   string `json:"passport_id"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
}

// ToDomain converts the HIS response to a domain model. ID, HospitalID and
// PatientID are left empty; the service fills them when it saves the patient.
// A record that cannot be stored — no national ID and no passport ID, an
// unparsable date of birth, an unknown gender, or no HN — yields the zero
// value.
func (m HospitalAPatient) ToDomain() domain.HospitalPatient {
	if m.NationalID == "" && m.PassportID == "" {
		return domain.HospitalPatient{}
	}

	dob, err := time.Parse(time.DateOnly, m.DateOfBirth)
	if err != nil {
		return domain.HospitalPatient{}
	}

	gender := enum.Gender(m.Gender)
	if !gender.IsValid() {
		return domain.HospitalPatient{}
	}

	if m.PatientHN == "" {
		return domain.HospitalPatient{}
	}

	return domain.HospitalPatient{
		PatientHN:   m.PatientHN,
		PhoneNumber: m.PhoneNumber,
		Email:       m.Email,
		Patient: domain.Patient{
			NationalID:   m.NationalID,
			PassportID:   m.PassportID,
			FirstNameTH:  m.FirstNameTH,
			MiddleNameTH: m.MiddleNameTH,
			LastNameTH:   m.LastNameTH,
			FirstNameEN:  m.FirstNameEN,
			MiddleNameEN: m.MiddleNameEN,
			LastNameEN:   m.LastNameEN,
			DateOfBirth:  dob,
			Gender:       gender,
		},
	}
}
