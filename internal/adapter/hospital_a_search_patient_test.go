package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/property"
)

const hospitalAPatientJSON = `{
	"first_name_th": "สมชาย", "middle_name_th": "", "last_name_th": "ใจดี",
	"first_name_en": "Somchai", "middle_name_en": "", "last_name_en": "Jaidee",
	"date_of_birth": "1990-05-17", "patient_hn": "HN0001",
	"national_id": "1234567890123", "passport_id": "",
	"phone_number": "0812345678", "email": "somchai@example.com", "gender": "M"
}`

func TestHospitalASearchPatient(t *testing.T) {
	ctx := context.Background()

	newAdapter := func(t *testing.T, h http.HandlerFunc, timeout time.Duration) *adapter {
		srv := httptest.NewServer(h)
		t.Cleanup(srv.Close)
		cfg := &property.Config{HISHospitalA: property.HISHospitalAConfig{
			HISHospitalASearchPatientURL: srv.URL + "/patient/search",
			HISHospitalATimeout:          timeout,
		}}
		return New(cfg).(*adapter)
	}
	respond := func(status int, body string) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
			_, _ = w.Write([]byte(body))
		}
	}

	t.Run("success", func(t *testing.T) {
		var path string
		a := newAdapter(t, func(w http.ResponseWriter, r *http.Request) {
			path = r.URL.EscapedPath()
			respond(http.StatusOK, hospitalAPatientJSON)(w, r)
		}, time.Second)

		hp, err := a.HospitalASearchPatient(ctx, "1234567890123")
		require.NoError(t, err)
		assert.Equal(t, "/patient/search/1234567890123", path)
		assert.Equal(t, "HN0001", hp.PatientHN)
		assert.Equal(t, "Somchai", hp.Patient.FirstNameEN)
	})

	t.Run("id stays a single path segment", func(t *testing.T) {
		var path string
		a := newAdapter(t, func(w http.ResponseWriter, r *http.Request) {
			path = r.URL.EscapedPath()
			respond(http.StatusNotFound, "")(w, r)
		}, time.Second)

		_, _ = a.HospitalASearchPatient(ctx, "../admin")
		assert.Equal(t, "/patient/search/..%2Fadmin", path)
	})

	tests := []struct {
		name    string
		handler http.HandlerFunc
		want    error
	}{
		{"not found", respond(http.StatusNotFound, ""), domain.ErrHISPatientNotFound},
		{"server error", respond(http.StatusInternalServerError, ""), domain.ErrHISUnavailable},
		{"malformed body", respond(http.StatusOK, "{"), domain.ErrHISInvalidResponse},
		{"record that can't be stored", respond(http.StatusOK, `{"patient_hn": "HN0001"}`), domain.ErrHISInvalidResponse},
		{"timeout", func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
		}, domain.ErrHISTimeout},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newAdapter(t, tt.handler, 50*time.Millisecond)
			_, err := a.HospitalASearchPatient(ctx, "1234567890123")
			assert.ErrorIs(t, err, tt.want)
		})
	}

	t.Run("unreachable", func(t *testing.T) {
		a := newAdapter(t, respond(http.StatusOK, ""), time.Second)
		a.cfg.HISHospitalA.HISHospitalASearchPatientURL = "http://127.0.0.1:1/patient/search"

		_, err := a.HospitalASearchPatient(ctx, "1234567890123")
		assert.ErrorIs(t, err, domain.ErrHISUnavailable)
	})
}
