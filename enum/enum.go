// Package enum holds the value sets shared by every layer (DTO, domain and
// entity), so none of them has to import another just to name a constant.
package enum

// Gender is a patient's gender as recorded by the hospital information system.
// It matches the gender column's CHECK constraint.
type Gender string

const (
	GenderMale   Gender = "M"
	GenderFemale Gender = "F"
)

// IsValid reports whether g is one of the supported genders.
func (g Gender) IsValid() bool {
	return g == GenderMale || g == GenderFemale
}

// HospitalCode is a hospitals.code value. Only hospitals whose HIS the
// middleware integrates with need a constant here.
type HospitalCode string

const (
	HospitalCodeA HospitalCode = "hospital-a"
)
