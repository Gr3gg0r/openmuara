package version

import (
	"strings"
	"testing"
)

func TestStringIncludesVersionCommitAndBuildTime(t *testing.T) {
	Version = "v1.2.3"
	Commit = "abc123"
	BuildTime = "2026-06-12T00:00:00Z"

	s := String()
	if !strings.Contains(s, Version) {
		t.Errorf("expected %q to contain version %q", s, Version)
	}
	if !strings.Contains(s, Commit) {
		t.Errorf("expected %q to contain commit %q", s, Commit)
	}
	if !strings.Contains(s, BuildTime) {
		t.Errorf("expected %q to contain build time %q", s, BuildTime)
	}
}
