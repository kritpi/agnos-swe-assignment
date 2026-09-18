package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/repository/entity"
)

func (r *Repository) UpsertHospitalPatient(ctx context.Context, hp domain.HospitalPatient) (*domain.HospitalPatient, error) {
	e := entity.HospitalPatientFromDomain(hp)
	sql := fmt.Sprintf(
		`INSERT INTO %s (id, hospital_id, patient_id, patient_hn, phone_number, email, created_at, updated_at)
		VALUES (@id, @hospital_id, @patient_id, @patient_hn, @phone_number, @email, @created_at, @updated_at)
		ON CONFLICT (hospital_id, patient_id) DO UPDATE SET
			patient_hn = EXCLUDED.patient_hn,
			phone_number = EXCLUDED.phone_number,
			email = EXCLUDED.email,
			updated_at = EXCLUDED.updated_at
		RETURNING id, hospital_id, patient_id, patient_hn, phone_number, email, created_at, updated_at`,
		table(r.tables.HospitalPatients))
	rows, err := r.conn(ctx).Query(ctx, sql, pgx.NamedArgs{
		"id":           e.ID,
		"hospital_id":  e.HospitalID,
		"patient_id":   e.PatientID,
		"patient_hn":   e.PatientHN,
		"phone_number": e.PhoneNumber,
		"email":        e.Email,
		"created_at":   e.CreatedAt,
		"updated_at":   e.UpdatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert hospital patient: %w", err)
	}

	stored, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.HospitalPatient])
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return nil, domain.ErrPatientConflict
	}
	if err != nil {
		return nil, fmt.Errorf("upsert hospital patient: %w", err)
	}

	d := stored.ToDomain()
	return &d, nil
}
