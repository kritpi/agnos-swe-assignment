package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/internal/core/port/mocks"
)

func TestVerifyStaff(t *testing.T) {
	ctx := context.Background()

	t.Run("member", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindStaffByID(ctx, "hosp-1", "staff-1").Return(&domain.Staff{ID: "staff-1"}, nil)

		assert.NoError(t, New(repo, nil, nil).VerifyStaff(ctx, "hosp-1", "staff-1"))
	})

	t.Run("not a member", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindStaffByID(ctx, "hosp-1", "staff-1").Return(nil, domain.ErrNotFound)

		assert.ErrorIs(t, New(repo, nil, nil).VerifyStaff(ctx, "hosp-1", "staff-1"), domain.ErrStaffNotMember)
	})

	t.Run("repository error", func(t *testing.T) {
		dbErr := errors.New("connection refused")
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindStaffByID(ctx, "hosp-1", "staff-1").Return(nil, dbErr)

		assert.ErrorIs(t, New(repo, nil, nil).VerifyStaff(ctx, "hosp-1", "staff-1"), dbErr)
	})
}
