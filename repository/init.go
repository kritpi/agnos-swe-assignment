package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kritpi/agnos-swe-assignment/internal/core/port"
)

// Repository is the concrete data-access implementation. The struct and its
// constructor live here; each operation is implemented as a receiver method
// in its own file (one file per operation).
type Repository struct {
	db *pgxpool.Pool
}

// New builds a Repository from a PostgreSQL connection pool.
func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

// DBTX is the query surface shared by *pgxpool.Pool and pgx.Tx.
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type txKey struct{}

// conn returns the transaction stored in ctx by Transactional, or the pool
// when called outside a transaction. Operation files should query through it.
func (r *Repository) conn(ctx context.Context) DBTX {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return r.db
}

// Transactional runs f inside a database transaction. It commits when f
// returns nil and rolls back on error or panic. Nested calls reuse the
// outer transaction.
func (r *Repository) Transactional(ctx context.Context, f func(ctx context.Context) error) (err error) {
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return f(ctx)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
				err = errors.Join(err, fmt.Errorf("rollback transaction: %w", rbErr))
			}
		}
	}()

	if err = f(context.WithValue(ctx, txKey{}, tx)); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// compile-time assertion that Repository satisfies the port.Repository interface.
var _ port.Repository = (*Repository)(nil)
