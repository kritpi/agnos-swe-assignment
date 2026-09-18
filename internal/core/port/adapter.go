package port

import (
	"context"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

// Adapter talks to the hospital information systems (HIS). It has one method
// per hospital operation, because each HIS has its own contract.
type Adapter interface {
	// HospitalASearchPatient looks a patient up in hospital A's HIS by
	// national ID or passport ID.
	HospitalASearchPatient(ctx context.Context, id string) (*domain.HospitalPatient, error)
}
