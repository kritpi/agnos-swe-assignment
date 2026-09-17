package port

import "context"

// Repository is the data-access interface. Each method is implemented in its
// own file under repository/.
type Repository interface {
	// Transactional runs f inside a database transaction. Repository calls made
	// with the ctx passed to f join that transaction.
	Transactional(ctx context.Context, f func(ctx context.Context) error) error
}
