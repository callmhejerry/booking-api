package config

import (
	"net/url"
	"testing"
)

func TestDSNUsesRenderDatabaseURLAndAddsSSLMode(t *testing.T) {
	cfg := dbConfig{
		DatabaseURL: "postgresql://booking:secret@db.example:5432/booking",
		SSLMode:     "require",
	}

	parsed, err := url.Parse(cfg.DSN())
	if err != nil {
		t.Fatalf("parse DSN: %v", err)
	}
	if parsed.Scheme != "postgresql" {
		t.Errorf("scheme = %q, want postgresql", parsed.Scheme)
	}
	if got := parsed.Query().Get("sslmode"); got != "require" {
		t.Errorf("sslmode = %q, want require", got)
	}
}

func TestDSNEncodesCredentialsFromIndividualSettings(t *testing.T) {
	cfg := dbConfig{
		Host:     "localhost",
		Port:     "5432",
		Name:     "booking",
		Username: "booking-user",
		Password: "pass@word",
		SSLMode:  "disable",
	}

	parsed, err := url.Parse(cfg.DSN())
	if err != nil {
		t.Fatalf("parse DSN: %v", err)
	}
	if got := parsed.User.Username(); got != "booking-user" {
		t.Errorf("username = %q, want booking-user", got)
	}
	password, ok := parsed.User.Password()
	if !ok || password != "pass@word" {
		t.Errorf("password = %q, want pass@word", password)
	}
	if got := parsed.Query().Get("sslmode"); got != "disable" {
		t.Errorf("sslmode = %q, want disable", got)
	}
}

func TestMigrationDSNUsesPGX5Scheme(t *testing.T) {
	cfg := dbConfig{DatabaseURL: "postgres://booking:secret@db.example:5432/booking"}

	parsed, err := url.Parse(cfg.MigrationDSN())
	if err != nil {
		t.Fatalf("parse migration DSN: %v", err)
	}
	if parsed.Scheme != "pgx5" {
		t.Errorf("scheme = %q, want pgx5", parsed.Scheme)
	}
}

func TestLoadConfigRequiresTLSForReleaseMode(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://booking:secret@db.example:5432/booking")
	t.Setenv("POSTGRES_SSL_MODE", "")
	t.Setenv("PORT", "10000")
	t.Setenv("MODE", "release")
	t.Setenv("API_KEY", "test-key")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if got := cfg.DBConfig.DSN(); !containsQueryValue(t, got, "sslmode", "require") {
		t.Errorf("DSN %q does not require TLS", got)
	}
}

func TestLoadConfigRejectsDatabaseURLWithTLSDisabledInReleaseMode(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://booking:secret@db.example:5432/booking?sslmode=disable")
	t.Setenv("PORT", "10000")
	t.Setenv("MODE", "release")
	t.Setenv("API_KEY", "test-key")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected release mode to reject sslmode=disable")
	}
}

func containsQueryValue(t *testing.T, rawURL, key, want string) bool {
	t.Helper()
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse URL: %v", err)
	}
	return parsed.Query().Get(key) == want
}
