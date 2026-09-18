package entity

import (
	"time"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

type Hospital struct {
	ID        string    `db:"id"`
	Code      string    `db:"code"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (e Hospital) ToDomain() domain.Hospital {
	return domain.Hospital{
		ID:        e.ID,
		Code:      e.Code,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func HospitalFromDomain(d domain.Hospital) Hospital {
	return Hospital{
		ID:        d.ID,
		Code:      d.Code,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
