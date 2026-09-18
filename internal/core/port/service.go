package port

import (
	"context"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

type Service interface {
	CreateStaff(ctx context.Context, cred domain.StaffCredentials) (*domain.Staff, error)
}
