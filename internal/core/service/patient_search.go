package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kritpi/agnos-swe-assignment/enum"
	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

func (s *service) SearchPatient(ctx context.Context, hospitalCode string, in domain.PatientSearchCriteria) (*domain.HospitalPatient, error) {
	if in.IsEmpty() {
		return nil, domain.ErrEmptySearchCriteria
	}

	// 1. search the local database
	found, err := s.repo.SearchPatients(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("search patients: %w", err)
	}
	switch {
	case len(found) == 1:
		return &found[0], nil
	case len(found) > 1:
		return nil, domain.ErrPatientAmbiguous
	case !in.HasPatientID():
		return nil, domain.ErrPatientNotFound
	}

	// 2. fall back to the staff member's own hospital HIS
	search, ok := s.hisSearch[enum.HospitalCode(hospitalCode)]
	if !ok {
		return nil, domain.ErrPatientNotFound
	}
	his, err := fetchFromHIS(ctx, search, in)
	if errors.Is(err, domain.ErrHISPatientNotFound) {
		return nil, domain.ErrPatientNotFound
	}
	if err != nil {
		return nil, err
	}

	// 3. store it, even if it fails the other filters: the record itself is valid
	var saved *domain.HospitalPatient
	save := func(ctx context.Context) error {
		now := time.Now().UTC()

		existing, err := s.findExistingPatient(ctx, his.Patient)
		if err != nil {
			return err
		}
		p, err := s.storePatient(ctx, existing, his.Patient, now)
		if err != nil {
			return err
		}
		saved, err = s.registerPatient(ctx, in.HospitalID, p, *his, now)
		return err
	}

	// A concurrent request may insert the same patient first. The retry then
	// finds that patient and updates it instead.
	err = s.repo.Transactional(ctx, save)
	if errors.Is(err, domain.ErrPatientAlreadyExists) {
		err = s.repo.Transactional(ctx, save)
	}
	if errors.Is(err, domain.ErrPatientConflict) || errors.Is(err, domain.ErrPatientAlreadyExists) {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("save patient: %w", err)
	}

	// 4. return it only if it matches every filter
	if !in.Matches(*saved) {
		return nil, domain.ErrPatientNotFound
	}
	return saved, nil
}

func fetchFromHIS(ctx context.Context, search hisSearchFunc, in domain.PatientSearchCriteria) (*domain.HospitalPatient, error) {
	for _, id := range []string{in.NationalID, in.PassportID} {
		if id == "" {
			continue
		}
		hp, err := search(ctx, id)
		if errors.Is(err, domain.ErrHISPatientNotFound) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("search hospital information system: %w", err)
		}
		// guard against a HIS answering with someone else's record
		if hp.Patient.NationalID != id && hp.Patient.PassportID != id {
			return nil, fmt.Errorf("search hospital information system: %w: record doesn't carry the requested id", domain.ErrHISInvalidResponse)
		}
		return hp, nil
	}
	return nil, domain.ErrHISPatientNotFound
}

func (s *service) findExistingPatient(ctx context.Context, p domain.Patient) (*domain.Patient, error) {
	var existing *domain.Patient
	if p.NationalID != "" {
		found, err := s.repo.FindPatientByNationalID(ctx, p.NationalID)
		switch {
		case errors.Is(err, domain.ErrNotFound):
		case err != nil:
			return nil, fmt.Errorf("find patient by national id: %w", err)
		default:
			existing = found
		}
	}
	if p.PassportID != "" {
		found, err := s.repo.FindPatientByPassportID(ctx, p.PassportID)
		switch {
		case errors.Is(err, domain.ErrNotFound):
		case err != nil:
			return nil, fmt.Errorf("find patient by passport id: %w", err)
		case existing != nil && existing.ID != found.ID:
			return nil, domain.ErrPatientConflict
		default:
			existing = found
		}
	}
	return existing, nil
}

func (s *service) storePatient(ctx context.Context, existing *domain.Patient, p domain.Patient, now time.Time) (domain.Patient, error) {
	p.UpdatedAt = now
	if existing != nil {
		p.ID = existing.ID
		p.CreatedAt = existing.CreatedAt
		if err := s.repo.UpdatePatient(ctx, p); err != nil {
			return domain.Patient{}, fmt.Errorf("update patient: %w", err)
		}
		return p, nil
	}

	id, err := uuid.NewV7()
	if err != nil {
		return domain.Patient{}, fmt.Errorf("generate patient id: %w", err)
	}
	p.ID = id.String()
	p.CreatedAt = now
	if err := s.repo.CreatePatient(ctx, p); err != nil {
		return domain.Patient{}, fmt.Errorf("create patient: %w", err)
	}
	return p, nil
}

func (s *service) registerPatient(ctx context.Context, hospitalID string, p domain.Patient, his domain.HospitalPatient, now time.Time) (*domain.HospitalPatient, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate hospital patient id: %w", err)
	}
	stored, err := s.repo.UpsertHospitalPatient(ctx, domain.HospitalPatient{
		ID:          id.String(),
		HospitalID:  hospitalID,
		PatientID:   p.ID,
		PatientHN:   his.PatientHN,
		PhoneNumber: his.PhoneNumber,
		Email:       his.Email,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, fmt.Errorf("register patient at hospital: %w", err)
	}
	stored.Patient = p
	return stored, nil
}
