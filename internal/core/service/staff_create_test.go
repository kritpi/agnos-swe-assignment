package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/internal/core/port/mocks"
	"github.com/kritpi/agnos-swe-assignment/property"
)

func TestCreateStaff(t *testing.T) {
	ctx := context.Background()
	in := domain.StaffCredentials{Username: "nurse01", Password: "P@ssw0rd!", HospitalCode: "hospital-a"}
	hospital := &domain.Hospital{ID: "hosp-1", Code: "hospital-a"}

	t.Run("hospital not found", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindHospitalByCode(ctx, in.HospitalCode).Return(nil, domain.ErrNotFound)

		_, err := New(repo, nil, &property.Config{}).CreateStaff(ctx, in)
		assert.ErrorIs(t, err, domain.ErrHospitalNotFound)
	})

	t.Run("username already exists at hospital", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindHospitalByCode(ctx, in.HospitalCode).Return(hospital, nil)
		repo.EXPECT().FindStaffByUsername(ctx, hospital.ID, in.Username).Return(&domain.Staff{}, nil)

		_, err := New(repo, nil, &property.Config{}).CreateStaff(ctx, in)
		assert.ErrorIs(t, err, domain.ErrStaffAlreadyExists)
		repo.AssertNotCalled(t, "CreateStaff", mock.Anything, mock.Anything)
	})

	t.Run("unique violation on insert", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindHospitalByCode(ctx, in.HospitalCode).Return(hospital, nil)
		repo.EXPECT().FindStaffByUsername(ctx, hospital.ID, in.Username).Return(nil, domain.ErrNotFound)
		repo.EXPECT().CreateStaff(ctx, mock.Anything).Return(domain.ErrStaffAlreadyExists)

		_, err := New(repo, nil, &property.Config{}).CreateStaff(ctx, in)
		assert.ErrorIs(t, err, domain.ErrStaffAlreadyExists)
	})

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindHospitalByCode(ctx, in.HospitalCode).Return(hospital, nil)
		repo.EXPECT().FindStaffByUsername(ctx, hospital.ID, in.Username).Return(nil, domain.ErrNotFound)
		repo.EXPECT().CreateStaff(ctx, mock.Anything).Return(nil)

		staff, err := New(repo, nil, &property.Config{}).CreateStaff(ctx, in)
		require.NoError(t, err)

		assert.Equal(t, hospital.ID, staff.HospitalID)
		assert.Equal(t, in.Username, staff.Username)
		id, err := uuid.Parse(staff.ID)
		require.NoError(t, err)
		assert.Equal(t, uuid.Version(7), id.Version())
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(in.Password)))
	})
}
