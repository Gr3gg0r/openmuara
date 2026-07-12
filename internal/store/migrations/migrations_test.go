package migrations

import (
	"strings"
	"testing"
)

func TestReadExistingMigration(t *testing.T) {
	got, err := Read("003_audit_logs.sql")
	if err != nil {
		t.Fatalf("Read existing migration: %v", err)
	}
	if !strings.Contains(got, "audit_logs") {
		t.Errorf("expected audit_logs schema, got:\n%s", got)
	}
}

func TestReadMissingMigration(t *testing.T) {
	_, err := Read("does_not_exist.sql")
	if err == nil {
		t.Fatal("expected error for missing migration")
	}
	if !strings.Contains(err.Error(), "does_not_exist.sql") {
		t.Errorf("error should mention file name: %v", err)
	}
}
