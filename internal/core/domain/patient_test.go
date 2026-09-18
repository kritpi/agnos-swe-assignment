package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPatientSearchCriteriaMatches(t *testing.T) {
	dob := time.Date(1990, 5, 17, 0, 0, 0, 0, time.UTC)
	otherDOB := time.Date(1991, 5, 17, 0, 0, 0, 0, time.UTC)
	hp := HospitalPatient{
		PhoneNumber: "0812345678",
		Email:       "Somchai@Example.com",
		Patient: Patient{
			NationalID:  "1234567890123",
			FirstNameTH: "สมชาย",
			FirstNameEN: "Somchai",
			LastNameEN:  "Jaidee",
			DateOfBirth: dob,
		},
	}

	tests := []struct {
		name string
		c    PatientSearchCriteria
		want bool
	}{
		{"no filters", PatientSearchCriteria{}, true},
		{"national id", PatientSearchCriteria{NationalID: "1234567890123"}, true},
		{"wrong national id", PatientSearchCriteria{NationalID: "0000000000000"}, false},
		{"passport id not on record", PatientSearchCriteria{PassportID: "AA1234567"}, false},
		{"thai first name", PatientSearchCriteria{FirstName: "สมชาย"}, true},
		{"english first name any case", PatientSearchCriteria{FirstName: "somchai"}, true},
		{"wrong last name", PatientSearchCriteria{LastName: "Rakdee"}, false},
		{"middle name not on record", PatientSearchCriteria{MiddleName: "Mid"}, false},
		{"date of birth", PatientSearchCriteria{DateOfBirth: &dob}, true},
		{"wrong date of birth", PatientSearchCriteria{DateOfBirth: &otherDOB}, false},
		{"phone", PatientSearchCriteria{PhoneNumber: "0812345678"}, true},
		{"wrong phone", PatientSearchCriteria{PhoneNumber: "0899999999"}, false},
		{"email any case", PatientSearchCriteria{Email: "somchai@example.com"}, true},
		{"all filters", PatientSearchCriteria{NationalID: "1234567890123", FirstName: "Somchai", LastName: "jaidee", DateOfBirth: &dob}, true},
		{"one filter off", PatientSearchCriteria{NationalID: "1234567890123", FirstName: "Somsak"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.Matches(hp))
		})
	}
}

func TestPatientSearchCriteriaIsEmpty(t *testing.T) {
	dob := time.Date(1990, 5, 17, 0, 0, 0, 0, time.UTC)
	assert.True(t, PatientSearchCriteria{}.IsEmpty())
	assert.True(t, PatientSearchCriteria{HospitalID: "hosp-a"}.IsEmpty(), "the hospital is not a filter")
	for _, c := range []PatientSearchCriteria{
		{NationalID: "1234567890123"},
		{PassportID: "AA1234567"},
		{FirstName: "Somchai"},
		{MiddleName: "Mid"},
		{LastName: "Jaidee"},
		{DateOfBirth: &dob},
		{PhoneNumber: "0812345678"},
		{Email: "somchai@example.com"},
	} {
		assert.False(t, c.IsEmpty(), "%+v", c)
	}
}
