package service

import (
	"context"

	"github.com/kritpi/agnos-swe-assignment/enum"
	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/internal/core/port"
	"github.com/kritpi/agnos-swe-assignment/property"
)

// hisSearchFunc looks a patient up in one hospital's HIS by national ID or
// passport ID.
type hisSearchFunc func(ctx context.Context, id string) (*domain.HospitalPatient, error)

type service struct {
	repo     port.Repository
	adapters port.Adapter
	cfg      *property.Config
	// hisSearch routes a hospital code to its HIS patient search. A hospital
	// missing from the map has no HIS integration.
	hisSearch map[enum.HospitalCode]hisSearchFunc
}

func New(repo port.Repository, adapters port.Adapter, cfg *property.Config) port.Service {
	hisSearch := map[enum.HospitalCode]hisSearchFunc{}
	if adapters != nil {
		hisSearch[enum.HospitalCodeA] = adapters.HospitalASearchPatient
	}

	return &service{
		repo:      repo,
		adapters:  adapters,
		cfg:       cfg,
		hisSearch: hisSearch,
	}
}
