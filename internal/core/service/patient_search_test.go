package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kritpi/agnos-swe-assignment/enum"
	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/internal/core/port/mocks"
)

func TestSearchPatient(t *testing.T) {
	ctx := context.Background()
	const hospitalID = "hosp-a"
	hospitalCode := string(enum.HospitalCodeA)
	dob := time.Date(1990, 5, 17, 0, 0, 0, 0, time.UTC)

	// his is a record as the adapter returns it: no IDs, no hospital.
	his := func() *domain.HospitalPatient {
		return &domain.HospitalPatient{
			PatientHN:   "HN0001",
			PhoneNumber: "0812345678",
			Email:       "somchai@example.com",
			Patient: domain.Patient{
				NationalID:  "1234567890123",
				FirstNameTH: "สมชาย",
				FirstNameEN: "Somchai",
				LastNameEN:  "Jaidee",
				DateOfBirth: dob,
				Gender:      enum.GenderMale,
			},
		}
	}
	byNationalID := domain.PatientSearchCriteria{HospitalID: hospitalID, NationalID: "1234567890123"}

	// runTx makes the Transactional mock run f, like the real repository.
	runTx := func(repo *mocks.Repository) {
		repo.EXPECT().Transactional(mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
	}
	// echoUpsert makes UpsertHospitalPatient return what it was given.
	echoUpsert := func(repo *mocks.Repository) {
		repo.EXPECT().UpsertHospitalPatient(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, hp domain.HospitalPatient) (*domain.HospitalPatient, error) { return &hp, nil })
	}

	t.Run("empty criteria are rejected", func(t *testing.T) {
		_, err := New(mocks.NewRepository(t), mocks.NewAdapter(t), nil).
			SearchPatient(ctx, hospitalCode, domain.PatientSearchCriteria{HospitalID: hospitalID})
		assert.ErrorIs(t, err, domain.ErrEmptySearchCriteria)
	})

	t.Run("local hit does not call the HIS", func(t *testing.T) {
		local := []domain.HospitalPatient{{PatientHN: "HN0001"}}
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byNationalID).Return(local, nil)
		adapters := mocks.NewAdapter(t)

		res, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, byNationalID)
		require.NoError(t, err)
		assert.Equal(t, &local[0], res)
	})

	t.Run("more than one local match is ambiguous", func(t *testing.T) {
		byName := domain.PatientSearchCriteria{HospitalID: hospitalID, FirstName: "Somchai"}
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byName).Return([]domain.HospitalPatient{{PatientHN: "HN0001"}, {PatientHN: "HN0002"}}, nil)

		_, err := New(repo, mocks.NewAdapter(t), nil).SearchPatient(ctx, hospitalCode, byName)
		assert.ErrorIs(t, err, domain.ErrPatientAmbiguous)
	})

	t.Run("local miss without an ID is not found", func(t *testing.T) {
		byName := domain.PatientSearchCriteria{HospitalID: hospitalID, FirstName: "Somchai"}
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byName).Return([]domain.HospitalPatient{}, nil)

		_, err := New(repo, mocks.NewAdapter(t), nil).SearchPatient(ctx, hospitalCode, byName)
		assert.ErrorIs(t, err, domain.ErrPatientNotFound)
	})

	t.Run("hospital without a HIS adapter is not found", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byNationalID).Return([]domain.HospitalPatient{}, nil)

		_, err := New(repo, mocks.NewAdapter(t), nil).SearchPatient(ctx, "hospital-z", byNationalID)
		assert.ErrorIs(t, err, domain.ErrPatientNotFound)
	})

	t.Run("HIS not found is not found", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byNationalID).Return([]domain.HospitalPatient{}, nil)
		adapters := mocks.NewAdapter(t)
		adapters.EXPECT().HospitalASearchPatient(ctx, "1234567890123").Return(nil, domain.ErrHISPatientNotFound)

		_, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, byNationalID)
		assert.ErrorIs(t, err, domain.ErrPatientNotFound)
	})

	for _, hisErr := range []error{domain.ErrHISUnavailable, domain.ErrHISTimeout, domain.ErrHISInvalidResponse} {
		t.Run("HIS error "+hisErr.Error(), func(t *testing.T) {
			repo := mocks.NewRepository(t)
			repo.EXPECT().SearchPatients(ctx, byNationalID).Return([]domain.HospitalPatient{}, nil)
			adapters := mocks.NewAdapter(t)
			adapters.EXPECT().HospitalASearchPatient(ctx, "1234567890123").Return(nil, hisErr)

			_, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, byNationalID)
			assert.ErrorIs(t, err, hisErr)
		})
	}

	t.Run("HIS record for a different ID is invalid", func(t *testing.T) {
		other := his()
		other.Patient.NationalID = "9999999999999"
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byNationalID).Return([]domain.HospitalPatient{}, nil)
		adapters := mocks.NewAdapter(t)
		adapters.EXPECT().HospitalASearchPatient(ctx, "1234567890123").Return(other, nil)

		_, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, byNationalID)
		assert.ErrorIs(t, err, domain.ErrHISInvalidResponse)
	})

	t.Run("new patient is fetched, stored and returned", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byNationalID).Return([]domain.HospitalPatient{}, nil)
		runTx(repo)
		repo.EXPECT().FindPatientByNationalID(mock.Anything, "1234567890123").Return(nil, domain.ErrNotFound)
		repo.EXPECT().CreatePatient(mock.Anything, mock.MatchedBy(func(p domain.Patient) bool {
			return p.ID != "" && p.NationalID == "1234567890123"
		})).Return(nil)
		echoUpsert(repo)
		adapters := mocks.NewAdapter(t)
		adapters.EXPECT().HospitalASearchPatient(ctx, "1234567890123").Return(his(), nil)

		res, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, byNationalID)
		require.NoError(t, err)
		assert.Equal(t, hospitalID, res.HospitalID)
		assert.Equal(t, "HN0001", res.PatientHN)
		assert.Equal(t, res.Patient.ID, res.PatientID)
	})

	t.Run("national ID not in HIS, passport retry succeeds", func(t *testing.T) {
		both := domain.PatientSearchCriteria{HospitalID: hospitalID, NationalID: "1234567890123", PassportID: "AA1234567"}
		byPassport := his()
		byPassport.Patient.NationalID = ""
		byPassport.Patient.PassportID = "AA1234567"

		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, both).Return([]domain.HospitalPatient{}, nil)
		runTx(repo)
		repo.EXPECT().FindPatientByPassportID(mock.Anything, "AA1234567").Return(nil, domain.ErrNotFound)
		repo.EXPECT().CreatePatient(mock.Anything, mock.Anything).Return(nil)
		echoUpsert(repo)
		adapters := mocks.NewAdapter(t)
		adapters.EXPECT().HospitalASearchPatient(ctx, "1234567890123").Return(nil, domain.ErrHISPatientNotFound)
		adapters.EXPECT().HospitalASearchPatient(ctx, "AA1234567").Return(byPassport, nil)

		// The request's national ID isn't on the record, so it doesn't match.
		_, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, both)
		assert.ErrorIs(t, err, domain.ErrPatientNotFound)
	})

	t.Run("record failing other filters is stored but not returned", func(t *testing.T) {
		wrongName := domain.PatientSearchCriteria{HospitalID: hospitalID, NationalID: "1234567890123", FirstName: "Somsak"}
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, wrongName).Return([]domain.HospitalPatient{}, nil)
		runTx(repo)
		repo.EXPECT().FindPatientByNationalID(mock.Anything, "1234567890123").Return(nil, domain.ErrNotFound)
		repo.EXPECT().CreatePatient(mock.Anything, mock.Anything).Return(nil)
		echoUpsert(repo)
		adapters := mocks.NewAdapter(t)
		adapters.EXPECT().HospitalASearchPatient(ctx, "1234567890123").Return(his(), nil)

		_, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, wrongName)
		assert.ErrorIs(t, err, domain.ErrPatientNotFound)
	})

	t.Run("patient from another hospital is replaced by the HIS record and linked", func(t *testing.T) {
		createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		existing := &domain.Patient{
			ID:          "patient-1",
			NationalID:  "1234567890123",
			PassportID:  "AA1234567",
			LastNameTH:  "ใจดี",
			DateOfBirth: dob,
			Gender:      enum.GenderMale,
			CreatedAt:   createdAt,
		}
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byNationalID).Return([]domain.HospitalPatient{}, nil)
		runTx(repo)
		repo.EXPECT().FindPatientByNationalID(mock.Anything, "1234567890123").Return(existing, nil)
		repo.EXPECT().UpdatePatient(mock.Anything, mock.MatchedBy(func(p domain.Patient) bool {
			// the HIS record wins, including the values it leaves empty
			return p.ID == "patient-1" && p.CreatedAt.Equal(createdAt) &&
				p.PassportID == "" && p.LastNameTH == "" && p.FirstNameEN == "Somchai"
		})).Return(nil)
		repo.EXPECT().UpsertHospitalPatient(mock.Anything, mock.MatchedBy(func(hp domain.HospitalPatient) bool {
			return hp.PatientID == "patient-1" && hp.HospitalID == hospitalID
		})).RunAndReturn(func(_ context.Context, hp domain.HospitalPatient) (*domain.HospitalPatient, error) { return &hp, nil })
		adapters := mocks.NewAdapter(t)
		adapters.EXPECT().HospitalASearchPatient(ctx, "1234567890123").Return(his(), nil)

		res, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, byNationalID)
		require.NoError(t, err)
		assert.Equal(t, "patient-1", res.Patient.ID)
	})

	t.Run("national and passport IDs of two different patients conflict", func(t *testing.T) {
		record := his()
		record.Patient.PassportID = "AA1234567"
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byNationalID).Return([]domain.HospitalPatient{}, nil)
		runTx(repo)
		repo.EXPECT().FindPatientByNationalID(mock.Anything, "1234567890123").Return(&domain.Patient{ID: "patient-1"}, nil)
		repo.EXPECT().FindPatientByPassportID(mock.Anything, "AA1234567").Return(&domain.Patient{ID: "patient-2"}, nil)
		adapters := mocks.NewAdapter(t)
		adapters.EXPECT().HospitalASearchPatient(ctx, "1234567890123").Return(record, nil)

		_, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, byNationalID)
		assert.ErrorIs(t, err, domain.ErrPatientConflict)
	})

	t.Run("stored patient with a contradicting ID is overwritten", func(t *testing.T) {
		record := his()
		record.Patient.PassportID = "AA1234567"
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byNationalID).Return([]domain.HospitalPatient{}, nil)
		runTx(repo)
		repo.EXPECT().FindPatientByNationalID(mock.Anything, "1234567890123").
			Return(&domain.Patient{ID: "patient-1", NationalID: "1234567890123", PassportID: "BB7654321"}, nil)
		repo.EXPECT().FindPatientByPassportID(mock.Anything, "AA1234567").Return(nil, domain.ErrNotFound)
		repo.EXPECT().UpdatePatient(mock.Anything, mock.MatchedBy(func(p domain.Patient) bool {
			return p.ID == "patient-1" && p.PassportID == "AA1234567"
		})).Return(nil)
		echoUpsert(repo)
		adapters := mocks.NewAdapter(t)
		adapters.EXPECT().HospitalASearchPatient(ctx, "1234567890123").Return(record, nil)

		res, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, byNationalID)
		require.NoError(t, err)
		assert.Equal(t, "AA1234567", res.Patient.PassportID)
	})

	t.Run("concurrent insert is retried as an update", func(t *testing.T) {
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byNationalID).Return([]domain.HospitalPatient{}, nil)
		runTx(repo)
		repo.EXPECT().FindPatientByNationalID(mock.Anything, "1234567890123").Return(nil, domain.ErrNotFound).Once()
		repo.EXPECT().CreatePatient(mock.Anything, mock.Anything).Return(domain.ErrPatientAlreadyExists).Once()
		repo.EXPECT().FindPatientByNationalID(mock.Anything, "1234567890123").
			Return(&domain.Patient{ID: "patient-1", NationalID: "1234567890123"}, nil).Once()
		repo.EXPECT().UpdatePatient(mock.Anything, mock.Anything).Return(nil)
		echoUpsert(repo)
		adapters := mocks.NewAdapter(t)
		adapters.EXPECT().HospitalASearchPatient(ctx, "1234567890123").Return(his(), nil)

		res, err := New(repo, adapters, nil).SearchPatient(ctx, hospitalCode, byNationalID)
		require.NoError(t, err)
		assert.Equal(t, "patient-1", res.PatientID)
	})

	t.Run("repository error", func(t *testing.T) {
		dbErr := errors.New("connection refused")
		repo := mocks.NewRepository(t)
		repo.EXPECT().SearchPatients(ctx, byNationalID).Return(nil, dbErr)

		_, err := New(repo, mocks.NewAdapter(t), nil).SearchPatient(ctx, hospitalCode, byNationalID)
		assert.ErrorIs(t, err, dbErr)
	})
}
