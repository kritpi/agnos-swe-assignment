package entity

import (
	"time"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

type Staff struct {
	ID           string    `db:"id"`
	HospitalID   string    `db:"hospital_id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (e Staff) ToDomain() domain.Staff {
	return domain.Staff{
		ID:           e.ID,
		HospitalID:   e.HospitalID,
		Username:     e.Username,
		PasswordHash: e.PasswordHash,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func StaffFromDomain(d domain.Staff) Staff {
	return Staff{
		ID:           d.ID,
		HospitalID:   d.HospitalID,
		Username:     d.Username,
		PasswordHash: d.PasswordHash,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
