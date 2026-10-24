package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

type dbConfig struct {
	DatabaseURL string
	Port        string
	Name        string
	Username    string
	Password    string
	Host        string
	SSLMode     string
}

func (dbc *dbConfig) DSN() string {
	if dbc.DatabaseURL != "" {
		dsn, _ := url.Parse(dbc.DatabaseURL)
		query := dsn.Query()
		if query.Get("sslmode") == "" && dbc.SSLMode != "" {
			query.Set("sslmode", dbc.SSLMode)
		}
		dsn.RawQuery = query.Encode()
		return dsn.String()
	}

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbc.Username, dbc.Password),
		Host:   net.JoinHostPort(dbc.Host, dbc.Port),
		Path:   "/" + dbc.Name,
	}
	query := url.Values{}
	if dbc.SSLMode != "" {
		query.Set("sslmode", dbc.SSLMode)
	}
	dsn.RawQuery = query.Encode()
	return dsn.String()
}

func (dbc *dbConfig) MigrationDSN() string {
	dsn, _ := url.Parse(dbc.DSN())
	dsn.Scheme = "pgx5"
	return dsn.String()
}

type appConfig struct {
	Port               string
	Mode               string
	ApiKey             string
	CORSAllowedOrigins []string
}

type Config struct {
	AppCongfig appConfig
	DBConfig   dbConfig
}

func LoadConfig() (*Config, error) {
	dbConfig := dbConfig{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("POSTGRES_PORT"),
		Name:        os.Getenv("POSTGRES_DB"),
		Username:    os.Getenv("POSTGRES_USER"),
		Password:    os.Getenv("POSTGRES_PASSWORD"),
		Host:        os.Getenv("POSTGRES_HOST"),
		SSLMode:     os.Getenv("POSTGRES_SSL_MODE"),
	}

	appConfig := appConfig{
		Port:               os.Getenv("PORT"),
		Mode:               os.Getenv("MODE"),
		ApiKey:             os.Getenv("API_KEY"),
		CORSAllowedOrigins: parseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS")),
	}
	if len(appConfig.CORSAllowedOrigins) == 0 {
		appConfig.CORSAllowedOrigins = []string{"http://localhost:5000"}
	}
	if appConfig.Mode == "release" && dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "require"
	}

	cfg := &Config{
		DBConfig:   dbConfig,
		AppCongfig: appConfig,
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func parseAllowedOrigins(value string) []string {
	var origins []string
	for _, origin := range strings.Split(value, ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

func (cfg *Config) validate() error {
	if cfg.AppCongfig.ApiKey == "" {
		return fmt.Errorf("API_KEY is required")
	}
	if cfg.AppCongfig.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	if cfg.DBConfig.DatabaseURL != "" {
		dsn, err := url.Parse(cfg.DBConfig.DatabaseURL)
		if err != nil || (dsn.Scheme != "postgres" && dsn.Scheme != "postgresql") || dsn.Host == "" {
			return fmt.Errorf("DATABASE_URL must be a valid PostgreSQL connection URL")
		}
	} else if cfg.DBConfig.Password == "" ||
		cfg.DBConfig.Name == "" ||
		cfg.DBConfig.Username == "" ||
		cfg.DBConfig.Host == "" ||
		cfg.DBConfig.Port == "" {
		return fmt.Errorf("POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_HOST, and POSTGRES_PORT are required")
	}

	if cfg.AppCongfig.Mode == "release" {
		sslMode := cfg.DBConfig.SSLMode
		if cfg.DBConfig.DatabaseURL != "" {
			dsn, _ := url.Parse(cfg.DBConfig.DatabaseURL)
			if configuredSSLMode := dsn.Query().Get("sslmode"); configuredSSLMode != "" {
				sslMode = configuredSSLMode
			}
		}
		if sslMode != "require" && sslMode != "verify-ca" && sslMode != "verify-full" {
			return fmt.Errorf("release mode requires database TLS with sslmode=require, verify-ca, or verify-full")
		}
	}

	return nil
}
