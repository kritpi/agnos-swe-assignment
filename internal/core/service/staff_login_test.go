package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/internal/core/port/mocks"
	"github.com/kritpi/agnos-swe-assignment/property"
)

func TestLoginStaff(t *testing.T) {
	ctx := context.Background()
	cfg := &property.Config{JWT: property.JWTConfig{Secret: "test-secret", TTL: 24 * time.Hour}}
	in := domain.StaffCredentials{Username: "nurse01", Password: "P@ssw0rd!", HospitalCode: "hospital-a"}
	hospital := &domain.Hospital{ID: "hosp-1", Code: "hospital-a"}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.MinCost)
	require.NoError(t, err)
	staff := &domain.Staff{ID: "staff-1", HospitalID: hospital.ID, Username: in.Username, PasswordHash: string(hash)}

	t.Run("hospital not found", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindHospitalByCode(ctx, in.HospitalCode).Return(nil, domain.ErrNotFound)

		_, err := New(repo, nil, cfg).LoginStaff(ctx, in)
		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})

	t.Run("staff not found", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindHospitalByCode(ctx, in.HospitalCode).Return(hospital, nil)
		repo.EXPECT().FindStaffByUsername(ctx, hospital.ID, in.Username).Return(nil, domain.ErrNotFound)

		_, err := New(repo, nil, cfg).LoginStaff(ctx, in)
		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})

	t.Run("wrong password", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindHospitalByCode(ctx, in.HospitalCode).Return(hospital, nil)
		repo.EXPECT().FindStaffByUsername(ctx, hospital.ID, in.Username).Return(staff, nil)

		wrong := in
		wrong.Password = "not-the-password"
		_, err := New(repo, nil, cfg).LoginStaff(ctx, wrong)
		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})

	t.Run("repository error", func(t *testing.T) {
		dbErr := errors.New("connection refused")
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindHospitalByCode(ctx, in.HospitalCode).Return(nil, dbErr)

		_, err := New(repo, nil, cfg).LoginStaff(ctx, in)
		assert.ErrorIs(t, err, dbErr)
		assert.NotErrorIs(t, err, domain.ErrInvalidCredentials)
	})

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().FindHospitalByCode(ctx, in.HospitalCode).Return(hospital, nil)
		repo.EXPECT().FindStaffByUsername(ctx, hospital.ID, in.Username).Return(staff, nil)

		res, err := New(repo, nil, cfg).LoginStaff(ctx, in)
		require.NoError(t, err)

		var claims domain.StaffClaims
		token, err := jwt.ParseWithClaims(res.AccessToken, &claims, func(*jwt.Token) (any, error) {
			return []byte(cfg.JWT.Secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		require.NoError(t, err)
		require.True(t, token.Valid)

		assert.Equal(t, staff.ID, claims.StaffID)
		assert.Equal(t, staff.ID, claims.Subject)
		assert.Equal(t, staff.Username, claims.Username)
		assert.Equal(t, hospital.ID, claims.HospitalID)
		assert.Equal(t, hospital.Code, claims.HospitalCode)
		assert.Equal(t, 24*time.Hour, claims.ExpiresAt.Sub(claims.IssuedAt.Time))
	})
}
