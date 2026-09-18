package property

import (
	"net"
	"net/url"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ServerConfig struct {
	Port string `envconfig:"PORT" default:"8080"`
}

type PostgresConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	DBName   string `envconfig:"DB_NAME"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

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

type DBTableConfig struct {
	Hospitals        string `envconfig:"DB_TABLE_HOSPITALS" default:"hospitals"`
	Staffs           string `envconfig:"DB_TABLE_STAFFS" default:"staffs"`
	Patients         string `envconfig:"DB_TABLE_PATIENTS" default:"patients"`
	HospitalPatients string `envconfig:"DB_TABLE_HOSPITAL_PATIENTS" default:"hospital_patients"`
}

type HISHospitalAConfig struct {
	HISHospitalASearchPatientURL string        `envconfig:"HIS_HOSPITAL_A_SEARCH_PATIENT_URL"`
	HISHospitalATimeout          time.Duration `envconfig:"HIS_HOSPITAL_A_TIMEOUT" default:"10s"`
}

// JWTConfig holds the settings for signing staff access tokens.
type JWTConfig struct {
	Secret string        `envconfig:"JWT_SECRET" required:"true"`
	TTL    time.Duration `envconfig:"JWT_TTL" default:"24h"`
}

type Config struct {
	Server       ServerConfig
	Postgres     PostgresConfig
	DBTable      DBTableConfig
	JWT          JWTConfig
	HISHospitalA HISHospitalAConfig
}

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
	if err := envconfig.Process("", &cfg.JWT); err != nil {
		return nil, err
	}
	if err := envconfig.Process("", &cfg.HISHospitalA); err != nil {
		return nil, err
	}

	return &cfg, nil
}
