package persistence

import (
	"strings"
	"testing"
)

func TestDatabaseParams_NormalizedDriver_DefaultsToMySQL(t *testing.T) {
	p := DatabaseParams{}
	if got := p.NormalizedDriver(); got != "mysql" {
		t.Fatalf("expected mysql, got %q", got)
	}
}

func TestDatabaseParams_BuildDSN_PrefersExplicitDSN(t *testing.T) {
	p := DatabaseParams{Driver: "mysql", DSN: "  some-dsn  "}
	dsn, err := p.BuildDSN()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dsn != "some-dsn" {
		t.Fatalf("expected trimmed DSN, got %q", dsn)
	}
}

func TestDatabaseParams_BuildDSN_SqliteEmpty(t *testing.T) {
	p := DatabaseParams{Driver: "sqlite"}
	dsn, err := p.BuildDSN()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dsn != "" {
		t.Fatalf("expected empty DSN for sqlite (connector handles in-memory), got %q", dsn)
	}
}

func TestDatabaseParams_BuildDSN_MySQLStructured(t *testing.T) {
	p := DatabaseParams{
		Driver:   "mysql",
		Host:     "db",
		Port:     3307,
		Database: "store",
		User:     "root",
		Password: "secret",
	}
	dsn, err := p.BuildDSN()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(dsn, "tcp(db:3307)") {
		t.Fatalf("expected tcp host:port in DSN, got %q", dsn)
	}
	if !strings.Contains(dsn, "/store") {
		t.Fatalf("expected db name in DSN, got %q", dsn)
	}
	if !strings.Contains(dsn, "charset=utf8mb4") {
		t.Fatalf("expected charset=utf8mb4 default, got %q", dsn)
	}
	// parseTime is a driver option; some DSN formatters omit it when it matches the default.
}

func TestDatabaseParams_BuildDSN_MySQLMissingUserOrDatabase(t *testing.T) {
	p := DatabaseParams{Driver: "mysql", Host: "db"}
	_, err := p.BuildDSN()
	if err == nil {
		t.Fatalf("expected error")
	}
}
