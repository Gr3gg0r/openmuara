package cli

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunDoctorAllToolsPresent(t *testing.T) {
	lookPath := func(string) (string, error) { return "/usr/bin/tool", nil }

	result := runDoctor(testDoctorCmd(), false, lookPath)

	if !result.Healthy {
		t.Fatalf("expected healthy result, got %+v", result)
	}
	for i, tool := range result.Tools {
		if !tool.Found {
			t.Errorf("expected %s to be found, got %+v", doctorTools[i].name, tool)
		}
	}
}

func TestRunDoctorMissingOptionalTool(t *testing.T) {
	lookPath := func(name string) (string, error) {
		if name == "govulncheck" {
			return "", fmt.Errorf("not found")
		}
		return "/usr/bin/tool", nil
	}

	result := runDoctor(testDoctorCmd(), false, lookPath)

	if !result.Healthy {
		t.Fatalf("optional tool missing should not make environment unhealthy: %+v", result)
	}
	govulncheckFound := false
	for _, tool := range result.Tools {
		if tool.Name == "govulncheck" {
			govulncheckFound = true
			if tool.Found {
				t.Error("expected govulncheck to be missing")
			}
		}
	}
	if !govulncheckFound {
		t.Error("expected govulncheck check in result")
	}
}

func TestRunDoctorMissingRequiredTool(t *testing.T) {
	lookPath := func(name string) (string, error) {
		if name == "go" {
			return "", fmt.Errorf("not found")
		}
		return "/usr/bin/tool", nil
	}

	result := runDoctor(testDoctorCmd(), false, lookPath)

	if result.Healthy {
		t.Fatal("expected unhealthy result for missing required tool")
	}
	goFound := false
	for _, tool := range result.Tools {
		if tool.Name == "go" {
			goFound = true
			if tool.Found {
				t.Error("expected go to be missing")
			}
		}
	}
	if !goFound {
		t.Error("expected go check in result")
	}
}

func TestNewDoctorCommand(t *testing.T) {
	cmd := newDoctorCommand()
	if cmd.Use != "doctor" {
		t.Fatalf("expected command use to be doctor, got %q", cmd.Use)
	}
}

func TestDoctorJSONSchema(t *testing.T) {
	lookPath := func(string) (string, error) { return "/usr/bin/tool", nil }
	jsonOutput = true
	defer func() { jsonOutput = false }()

	result := runDoctor(testDoctorCmd(), false, lookPath)

	var buf bytes.Buffer
	if err := respondJSON(&buf, result); err != nil {
		t.Fatalf("encode JSON: %v", err)
	}
	out := buf.String()
	for _, key := range []string{"\"healthy\"", "\"tools\"", "\"config\"", "\"providers\"", "\"webhook\"", "\"version\""} {
		if !strings.Contains(out, key) {
			t.Errorf("expected JSON to contain %s, got:\n%s", key, out)
		}
	}
}

func testDoctorCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	return cmd
}
