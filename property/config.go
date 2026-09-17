package property

import (
	"net"
	"net/url"

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

// Config is the aggregated application configuration.
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
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

	return &cfg, nil
}
