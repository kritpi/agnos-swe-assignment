package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kritpi/agnos-swe-assignment/property"
)

type Repository struct {
	db     *pgxpool.Pool
	tables property.DBTableConfig
}

func New(db *pgxpool.Pool, tables property.DBTableConfig) *Repository {
	return &Repository{
		db:     db,
		tables: tables,
	}
}

func table(name string) string {
	return pgx.Identifier{name}.Sanitize()
}

type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
