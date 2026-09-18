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

const pgUniqueViolation = "23505"

func (r *Repository) CreateStaff(ctx context.Context, s domain.Staff) error {
	e := entity.StaffFromDomain(s)
	sql := fmt.Sprintf(
		`INSERT INTO %s (id, hospital_id, username, password_hash, created_at, updated_at)
		VALUES (@id, @hospital_id, @username, @password_hash, @created_at, @updated_at)`,
		table(r.tables.Staffs))
	_, err := r.conn(ctx).Exec(ctx, sql, pgx.NamedArgs{
		"id":            e.ID,
		"hospital_id":   e.HospitalID,
		"username":      e.Username,
		"password_hash": e.PasswordHash,
		"created_at":    e.CreatedAt,
		"updated_at":    e.UpdatedAt,
	})

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return domain.ErrStaffAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("insert staff: %w", err)
	}
	return nil
}
