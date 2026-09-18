package adapter

import (
	"context"
	"encoding/json"
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
// which is either a national ID or a passport ID. Any failure is returned as
// is; only a successful response is decoded and converted.
func (a *adapter) HospitalASearchPatient(ctx context.Context, id string) (*domain.HospitalPatient, error) {
	cfg := a.cfg.HISHospitalA

	endpoint, err := url.JoinPath(cfg.HISHospitalASearchPatientURL, id)
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
		return nil, fmt.Errorf("call hospital A his: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call hospital A his: unexpected status %d", res.StatusCode)
	}

	var patient domainadapter.HospitalAPatient
	if err := json.NewDecoder(io.LimitReader(res.Body, maxResponseBytes)).Decode(&patient); err != nil {
		return nil, fmt.Errorf("decode hospital A his response: %w", err)
	}
	patientResp := patient.ToDomain()

	return &patientResp, nil
}
