package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/repository/entity"
)

func (r *Repository) FindPatientByNationalID(ctx context.Context, nationalID string) (*domain.Patient, error) {
	return r.findPatientBy(ctx, "national_id", nationalID)
}

func (r *Repository) FindPatientByPassportID(ctx context.Context, passportID string) (*domain.Patient, error) {
	return r.findPatientBy(ctx, "passport_id", passportID)
}

func (r *Repository) findPatientBy(ctx context.Context, column, value string) (*domain.Patient, error) {
	sql := fmt.Sprintf(
		`SELECT id, national_id, passport_id, first_name_th, first_name_en, middle_name_th, middle_name_en,
			last_name_th, last_name_en, date_of_birth, gender, created_at, updated_at
		FROM %s WHERE %s = @value
		FOR UPDATE`,
		table(r.tables.Patients), pgx.Identifier{column}.Sanitize())
	rows, err := r.conn(ctx).Query(ctx, sql, pgx.NamedArgs{"value": value})
	if err != nil {
		return nil, fmt.Errorf("query patient: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Patient])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan patient: %w", err)
	}

	p := e.ToDomain()
	return &p, nil
}
