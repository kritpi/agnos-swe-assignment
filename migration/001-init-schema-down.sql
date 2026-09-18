-- Drop tables in reverse dependency order (indexes are dropped with their tables)
DROP TABLE IF EXISTS hospital_patients;

DROP TABLE IF EXISTS patients;

DROP TABLE IF EXISTS staffs;

DROP TABLE IF EXISTS hospitals;
