-- Create table hospitals
CREATE TABLE
    hospitals (
        id TEXT PRIMARY KEY, -- UUID
        code TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

-- Create table staffs
CREATE TABLE
    staffs (
        id TEXT PRIMARY KEY, -- UUID
        hospital_id TEXT NOT NULL REFERENCES hospitals (id),
        username TEXT NOT NULL,
        password_hash TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        UNIQUE (hospital_id, username)
    );

-- Create table patients
CREATE TABLE
    patients (
        id TEXT PRIMARY KEY, -- UUID
        national_id TEXT UNIQUE,
        passport_id TEXT UNIQUE,
        first_name_th TEXT,
        first_name_en TEXT,
        middle_name_th TEXT,
        middle_name_en TEXT,
        last_name_th TEXT,
        last_name_en TEXT,
        date_of_birth DATE NOT NULL,
        gender TEXT NOT NULL CHECK (gender IN ('M', 'F')),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        CONSTRAINT chk_national_id_or_passport_id CHECK (
            national_id IS NOT NULL
            OR passport_id IS NOT NULL
        )
    );

-- Create table hospital_patients
CREATE TABLE
    hospital_patients (
        id TEXT PRIMARY KEY, -- UUID
        hospital_id TEXT NOT NULL REFERENCES hospitals (id),
        patient_id TEXT NOT NULL REFERENCES patients (id),
        patient_hn TEXT NOT NULL,
        phone_number TEXT NOT NULL,
        email TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        UNIQUE (hospital_id, patient_id),
        UNIQUE (hospital_id, patient_hn)
    );

-- Create index patients table
CREATE INDEX idx_patients_national_id ON patients (national_id);
CREATE INDEX idx_patients_passport_id ON patients (passport_id);