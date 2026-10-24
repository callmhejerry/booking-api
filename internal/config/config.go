package config

import (
	"fmt"
	"os"
)

type dbConfig struct {
	Port     string
	Name     string
	Username string
	Password string
	Host     string
	SSLMode  string
}

func (dbc *dbConfig) DSN() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", dbc.Username, dbc.Password, dbc.Host, dbc.Port, dbc.Name, dbc.SSLMode)
}

type appConfig struct {
	Port   string
	Mode   string
	ApiKey string
}

type Config struct {
	AppCongfig appConfig
	DBConfig   dbConfig
}

func LoadConfig() (*Config, error) {
	dbConfig := dbConfig{
		Port:     os.Getenv("POSTGRES_PORT"),
		Name:     os.Getenv("POSTGRES_DB"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		SSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
	}

	appConfig := appConfig{
		Port:   os.Getenv("PORT"),
		Mode:   os.Getenv("MODE"),
		ApiKey: os.Getenv("API_KEY"),
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

func (cfg *Config) validate() error {
	if cfg.AppCongfig.ApiKey == "" {
		return fmt.Errorf("Invalid api key")
	}
	if cfg.AppCongfig.Port == "" {
		return fmt.Errorf("invalid app port")
	}
	if cfg.DBConfig.Password == "" || cfg.DBConfig.Name == "" || cfg.DBConfig.Username == "" {
		return fmt.Errorf("Invalid database password")
	}

	if cfg.AppCongfig.Mode == "prod" && cfg.DBConfig.SSLMode == "disabled" {
		return fmt.Errorf("Cannot disable database ssl mode for production")
	}

	return nil
}
