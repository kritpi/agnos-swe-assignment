package property

import (
	"net"
	"net/url"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port string `envconfig:"PORT" default:"8080"`
}

// PostgresConfig holds the PostgreSQL connection settings.
type PostgresConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	DBName   string `envconfig:"DB_NAME"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

// DSN builds a postgres:// connection URL for pgx.
func (p PostgresConfig) DSN() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(p.User, p.Password),
		Host:     net.JoinHostPort(p.Host, p.Port),
		Path:     p.DBName,
		RawQuery: url.Values{"sslmode": []string{p.SSLMode}}.Encode(),
	}
	return u.String()
}

// DBTableConfig holds the database table names used by the repository. They
// must match the tables created by the migrations.
type DBTableConfig struct {
	Hospitals        string `envconfig:"DB_TABLE_HOSPITALS" default:"hospitals"`
	Staffs           string `envconfig:"DB_TABLE_STAFFS" default:"staffs"`
	Patients         string `envconfig:"DB_TABLE_PATIENTS" default:"patients"`
	HospitalPatients string `envconfig:"DB_TABLE_HOSPITAL_PATIENTS" default:"hospital_patients"`
}

// HISHospitalAConfig holds hospital A's information system settings. Each
// hospital gets its own struct, because each HIS has its own endpoints. URLs
// are full endpoint URLs; the adapter only appends the path parameter.
type HISHospitalAConfig struct {
	HISHospitalASearchPatientURL string        `envconfig:"HIS_HOSPITAL_A_SEARCH_PATIENT_URL"`
	HISHospitalATimeout          time.Duration `envconfig:"HIS_HOSPITAL_A_TIMEOUT" default:"10s"`
}

// Config is the aggregated application configuration.
type Config struct {
	Server       ServerConfig
	Postgres     PostgresConfig
	DBTable      DBTableConfig
	HISHospitalA HISHospitalAConfig
}

// Load reads the .env file (if present) and populates Config from the environment.
func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	var cfg Config
	if err := envconfig.Process("", &cfg.Server); err != nil {
		return nil, err
	}
	if err := envconfig.Process("", &cfg.Postgres); err != nil {
		return nil, err
	}
	if err := envconfig.Process("", &cfg.DBTable); err != nil {
		return nil, err
	}
	if err := envconfig.Process("", &cfg.HISHospitalA); err != nil {
		return nil, err
	}

	return &cfg, nil
}
