package entity

import (
	"time"

	"github.com/kritpi/agnos-swe-assignment/enum"
	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

// Patient maps to the patients table. Nullable columns are pointers.
type Patient struct {
	ID           string    `db:"id"`
	NationalID   *string   `db:"national_id"`
	PassportID   *string   `db:"passport_id"`
	FirstNameTH  *string   `db:"first_name_th"`
	FirstNameEN  *string   `db:"first_name_en"`
	MiddleNameTH *string   `db:"middle_name_th"`
	MiddleNameEN *string   `db:"middle_name_en"`
	LastNameTH   *string   `db:"last_name_th"`
	LastNameEN   *string   `db:"last_name_en"`
	DateOfBirth  time.Time `db:"date_of_birth"`
	Gender       string    `db:"gender"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// ToDomain converts the entity to its domain model.
func (e Patient) ToDomain() domain.Patient {
	return domain.Patient{
		ID:           e.ID,
		NationalID:   fromNullString(e.NationalID),
		PassportID:   fromNullString(e.PassportID),
		FirstNameTH:  fromNullString(e.FirstNameTH),
		MiddleNameTH: fromNullString(e.MiddleNameTH),
		LastNameTH:   fromNullString(e.LastNameTH),
		FirstNameEN:  fromNullString(e.FirstNameEN),
		MiddleNameEN: fromNullString(e.MiddleNameEN),
		LastNameEN:   fromNullString(e.LastNameEN),
		DateOfBirth:  e.DateOfBirth,
		Gender:       enum.Gender(e.Gender),
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

// PatientFromDomain converts a domain model to its entity.
func PatientFromDomain(d domain.Patient) Patient {
	return Patient{
		ID:           d.ID,
		NationalID:   toNullString(d.NationalID),
		PassportID:   toNullString(d.PassportID),
		FirstNameTH:  toNullString(d.FirstNameTH),
		MiddleNameTH: toNullString(d.MiddleNameTH),
		LastNameTH:   toNullString(d.LastNameTH),
		FirstNameEN:  toNullString(d.FirstNameEN),
		MiddleNameEN: toNullString(d.MiddleNameEN),
		LastNameEN:   toNullString(d.LastNameEN),
		DateOfBirth:  d.DateOfBirth,
		Gender:       string(d.Gender),
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
