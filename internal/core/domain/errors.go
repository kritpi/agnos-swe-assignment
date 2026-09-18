package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrHospitalNotFound   = errors.New("hospital not found")
	ErrStaffAlreadyExists = errors.New("username already exists at this hospital")
	ErrInvalidCredentials = errors.New("invalid username, password or hospital")
	ErrStaffNotMember     = errors.New("staff is not a member of this hospital")

	ErrPatientConflict      = errors.New("patient record conflicts with an existing patient")
	ErrPatientAlreadyExists = errors.New("patient already exists")
	ErrPatientNotFound      = errors.New("patient not found")
	ErrPatientAmbiguous     = errors.New("more than one patient matches, add more filters")
	ErrEmptySearchCriteria  = errors.New("at least one search filter is required")

	// HIS (hospital information system) errors, returned by the adapters.
	ErrHISPatientNotFound = errors.New("patient not found in hospital information system")
	ErrHISUnavailable     = errors.New("hospital information system is unavailable")
	ErrHISTimeout         = errors.New("hospital information system timed out")
	ErrHISInvalidResponse = errors.New("hospital information system returned an invalid response")
)
