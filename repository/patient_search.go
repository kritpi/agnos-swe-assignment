package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
	"github.com/kritpi/agnos-swe-assignment/repository/entity"
)

func (r *Repository) SearchPatients(ctx context.Context, c domain.PatientSearchCriteria) ([]domain.HospitalPatient, error) {
	where := []string{"hp.hospital_id = @hospital_id"}
	args := pgx.NamedArgs{"hospital_id": c.HospitalID}

	if c.NationalID != "" {
		where = append(where, "p.national_id = @national_id")
		args["national_id"] = c.NationalID
	}
	if c.PassportID != "" {
		where = append(where, "p.passport_id = @passport_id")
		args["passport_id"] = c.PassportID
	}
	if c.FirstName != "" {
		where = append(where, "(LOWER(p.first_name_th) = LOWER(@first_name) OR LOWER(p.first_name_en) = LOWER(@first_name))")
		args["first_name"] = c.FirstName
	}
	if c.MiddleName != "" {
		where = append(where, "(LOWER(p.middle_name_th) = LOWER(@middle_name) OR LOWER(p.middle_name_en) = LOWER(@middle_name))")
		args["middle_name"] = c.MiddleName
	}
	if c.LastName != "" {
		where = append(where, "(LOWER(p.last_name_th) = LOWER(@last_name) OR LOWER(p.last_name_en) = LOWER(@last_name))")
		args["last_name"] = c.LastName
	}
	if c.DateOfBirth != nil {
		where = append(where, "p.date_of_birth = @date_of_birth")
		args["date_of_birth"] = *c.DateOfBirth
	}
	if c.PhoneNumber != "" {
		where = append(where, "hp.phone_number = @phone_number")
		args["phone_number"] = c.PhoneNumber
	}
	if c.Email != "" {
		where = append(where, "LOWER(hp.email) = LOWER(@email)")
		args["email"] = c.Email
	}

	sql := fmt.Sprintf(
		`SELECT hp.id, hp.hospital_id, hp.patient_id, hp.patient_hn, hp.phone_number, hp.email, hp.created_at, hp.updated_at,
			p.id, p.national_id, p.passport_id, p.first_name_th, p.first_name_en, p.middle_name_th, p.middle_name_en,
			p.last_name_th, p.last_name_en, p.date_of_birth, p.gender, p.created_at, p.updated_at
		FROM %s hp JOIN %s p ON p.id = hp.patient_id
		WHERE %s
		ORDER BY hp.patient_hn
		LIMIT 2`,
		table(r.tables.HospitalPatients), table(r.tables.Patients), strings.Join(where, " AND "))
	rows, err := r.conn(ctx).Query(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("query patients: %w", err)
	}

	res, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.HospitalPatient, error) {
		var hp entity.HospitalPatient
		var p entity.Patient
		err := row.Scan(
			&hp.ID, &hp.HospitalID, &hp.PatientID, &hp.PatientHN, &hp.PhoneNumber, &hp.Email, &hp.CreatedAt, &hp.UpdatedAt,
			&p.ID, &p.NationalID, &p.PassportID, &p.FirstNameTH, &p.FirstNameEN, &p.MiddleNameTH, &p.MiddleNameEN,
			&p.LastNameTH, &p.LastNameEN, &p.DateOfBirth, &p.Gender, &p.CreatedAt, &p.UpdatedAt,
		)
		d := hp.ToDomain()
		d.Patient = p.ToDomain()
		return d, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan patients: %w", err)
	}
	if res == nil {
		res = []domain.HospitalPatient{}
	}

	return res, nil
}
