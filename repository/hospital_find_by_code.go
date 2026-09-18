package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/repository/entity"
)

func (r *Repository) FindHospitalByCode(ctx context.Context, code string) (*domain.Hospital, error) {
	sql := fmt.Sprintf(
		`SELECT id, code, name, created_at, updated_at FROM %s WHERE code = @code`,
		table(r.tables.Hospitals))
	rows, err := r.conn(ctx).Query(ctx, sql, pgx.NamedArgs{"code": code})
	if err != nil {
		return nil, fmt.Errorf("query hospital: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Hospital])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan hospital: %w", err)
	}

	h := e.ToDomain()
	return &h, nil
}
