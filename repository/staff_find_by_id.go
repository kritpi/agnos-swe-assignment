package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/repository/entity"
)

func (r *Repository) FindStaffByID(ctx context.Context, hospitalID, staffID string) (*domain.Staff, error) {
	sql := fmt.Sprintf(
		`SELECT id, hospital_id, username, password_hash, created_at, updated_at
		FROM %s WHERE hospital_id = @hospital_id AND id = @id`,
		table(r.tables.Staffs))
	rows, err := r.conn(ctx).Query(ctx, sql, pgx.NamedArgs{
		"hospital_id": hospitalID,
		"id":          staffID,
	})
	if err != nil {
		return nil, fmt.Errorf("query staff: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Staff])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan staff: %w", err)
	}

	s := e.ToDomain()
	return &s, nil
}
