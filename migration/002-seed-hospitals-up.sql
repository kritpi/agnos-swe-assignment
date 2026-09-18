-- Seed hospitals
INSERT INTO
    hospitals (id, code, name, created_at, updated_at)
VALUES
    (
        '01a0b3cd-23cd-73ea-9d59-5fdc5d0801f4',
        'hospital-a',
        'Hospital A',
        now (),
        now ()
    )
ON CONFLICT (id) DO NOTHING;
