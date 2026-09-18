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

func (r *Repository) CreatePatient(ctx context.Context, p domain.Patient) error {
	e := entity.PatientFromDomain(p)
	sql := fmt.Sprintf(
		`INSERT INTO %s (id, national_id, passport_id, first_name_th, first_name_en, middle_name_th, middle_name_en,
			last_name_th, last_name_en, date_of_birth, gender, created_at, updated_at)
		VALUES (@id, @national_id, @passport_id, @first_name_th, @first_name_en, @middle_name_th, @middle_name_en,
			@last_name_th, @last_name_en, @date_of_birth, @gender, @created_at, @updated_at)`,
		table(r.tables.Patients))
	_, err := r.conn(ctx).Exec(ctx, sql, patientArgs(e))

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return domain.ErrPatientAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("insert patient: %w", err)
	}
	return nil
}

// patientArgs binds every patients column. Shared by create and update.
func patientArgs(e entity.Patient) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":             e.ID,
		"national_id":    e.NationalID,
		"passport_id":    e.PassportID,
		"first_name_th":  e.FirstNameTH,
		"first_name_en":  e.FirstNameEN,
		"middle_name_th": e.MiddleNameTH,
		"middle_name_en": e.MiddleNameEN,
		"last_name_th":   e.LastNameTH,
		"last_name_en":   e.LastNameEN,
		"date_of_birth":  e.DateOfBirth,
		"gender":         e.Gender,
		"created_at":     e.CreatedAt,
		"updated_at":     e.UpdatedAt,
	}
}
