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
	Port string
	Mode string
}

type Config struct {
	AppCongfig appConfig
	DBConfig   dbConfig
}

func LoadConfig() *Config {
	dbConfig := dbConfig{
		Port:     os.Getenv("POSTGRES_PORT"),
		Name:     os.Getenv("POSTGRES_DB"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		SSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
	}

	appConfig := appConfig{
		Port: os.Getenv("PORT"),
	}

	return &Config{
		DBConfig:   dbConfig,
		AppCongfig: appConfig,
	}
}
