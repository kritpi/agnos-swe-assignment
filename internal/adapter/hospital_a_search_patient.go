package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	domainadapter "github.com/kritpi/agnos-swe-assignment/internal/adapter/domain-adapter"
	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

// maxResponseBytes caps how much of a HIS response is read.
const maxResponseBytes = 1 << 20 // 1 MB

// HospitalASearchPatient calls hospital A's patient search endpoint with id,
// which is either a national ID or a passport ID. Failures are classified as
// domain.ErrHISPatientNotFound (404), domain.ErrHISTimeout,
// domain.ErrHISUnavailable (network error or any other status) or
// domain.ErrHISInvalidResponse (a body that can't be decoded or stored).
func (a *adapter) HospitalASearchPatient(ctx context.Context, id string) (*domain.HospitalPatient, error) {
	cfg := a.cfg.HISHospitalA

	// Escape id so it stays a single path segment.
	endpoint, err := url.JoinPath(cfg.HISHospitalASearchPatientURL, url.PathEscape(id))
	if err != nil {
		return nil, fmt.Errorf("build hospital A search patient url: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.HISHospitalATimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build hospital A search patient request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := a.httpClient.Do(req)
	if err != nil {
		return nil, hisCallError("call hospital A his", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, domain.ErrHISPatientNotFound
	default:
		return nil, fmt.Errorf("call hospital A his: %w: unexpected status %d", domain.ErrHISUnavailable, res.StatusCode)
	}

	var patient domainadapter.HospitalAPatient
	if err := json.NewDecoder(io.LimitReader(res.Body, maxResponseBytes)).Decode(&patient); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, hisCallError("read hospital A his response", err)
		}
		return nil, fmt.Errorf("decode hospital A his response: %w: %w", domain.ErrHISInvalidResponse, err)
	}

	// ToDomain returns the zero value for a record that can't be stored.
	patientResp := patient.ToDomain()
	if patientResp.PatientHN == "" {
		return nil, fmt.Errorf("convert hospital A his response: %w", domain.ErrHISInvalidResponse)
	}

	return &patientResp, nil
}

// hisCallError classifies a transport error as a timeout or an unavailable HIS.
func hisCallError(op string, err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: %w: %w", op, domain.ErrHISTimeout, err)
	}
	return fmt.Errorf("%s: %w: %w", op, domain.ErrHISUnavailable, err)
}
