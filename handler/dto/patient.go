package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

type SearchPatientRequest struct {
	NationalID  string `json:"national_id" form:"national_id"`
	PassportID  string `json:"passport_id" form:"passport_id"`
	FirstName   string `json:"first_name" form:"first_name"`
	MiddleName  string `json:"middle_name" form:"middle_name"`
	LastName    string `json:"last_name" form:"last_name"`
	DateOfBirth string `json:"date_of_birth" form:"date_of_birth"` // YYYY-MM-DD
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Email       string `json:"email" form:"email"`
}

func (r SearchPatientRequest) ToDomain(hospitalID string) (domain.PatientSearchCriteria, error) {
	criteria := domain.PatientSearchCriteria{
		HospitalID:  hospitalID,
		NationalID:  strings.TrimSpace(r.NationalID),
		PassportID:  strings.TrimSpace(r.PassportID),
		FirstName:   strings.TrimSpace(r.FirstName),
		MiddleName:  strings.TrimSpace(r.MiddleName),
		LastName:    strings.TrimSpace(r.LastName),
		PhoneNumber: strings.TrimSpace(r.PhoneNumber),
		Email:       strings.TrimSpace(r.Email),
	}

	if dob := strings.TrimSpace(r.DateOfBirth); dob != "" {
		t, err := time.Parse(time.DateOnly, dob)
		if err != nil {
			return domain.PatientSearchCriteria{}, fmt.Errorf("invalid date_of_birth %q: expected YYYY-MM-DD", dob)
		}
		criteria.DateOfBirth = &t
	}

	return criteria, nil
}

type PatientResponse struct {
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

func NewPatientResponse(hp domain.HospitalPatient) PatientResponse {
	p := hp.Patient
	return PatientResponse{
		FirstNameTH:  p.FirstNameTH,
		MiddleNameTH: p.MiddleNameTH,
		LastNameTH:   p.LastNameTH,
		FirstNameEN:  p.FirstNameEN,
		MiddleNameEN: p.MiddleNameEN,
		LastNameEN:   p.LastNameEN,
		DateOfBirth:  p.DateOfBirth.Format(time.DateOnly),
		PatientHN:    hp.PatientHN,
		NationalID:   p.NationalID,
		PassportID:   p.PassportID,
		PhoneNumber:  hp.PhoneNumber,
		Email:        hp.Email,
		Gender:       string(p.Gender),
	}
}
