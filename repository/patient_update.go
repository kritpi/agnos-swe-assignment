package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/repository/entity"
)

func (r *Repository) UpdatePatient(ctx context.Context, p domain.Patient) error {
	e := entity.PatientFromDomain(p)
	sql := fmt.Sprintf(
		`UPDATE %s SET
			national_id = @national_id, passport_id = @passport_id,
			first_name_th = @first_name_th, first_name_en = @first_name_en,
			middle_name_th = @middle_name_th, middle_name_en = @middle_name_en,
			last_name_th = @last_name_th, last_name_en = @last_name_en,
			date_of_birth = @date_of_birth, gender = @gender, updated_at = @updated_at
		WHERE id = @id`,
		table(r.tables.Patients))
	tag, err := r.conn(ctx).Exec(ctx, sql, patientArgs(e))

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return domain.ErrPatientAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("update patient: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
