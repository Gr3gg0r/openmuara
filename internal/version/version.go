// Package version exposes build-time version metadata.
package version

import (
	"fmt"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"
)

var (
	// Version is the semantic version. Injected at build time.
	Version = "dev"
	// Commit is the git commit SHA. Injected at build time.
	Commit = "unknown"
	// BuildTime is the ISO8601 build timestamp. Injected at build time.
	BuildTime = "unknown"
)

func init() {
	// Fall back to Go build metadata when ldflags are not injected. This keeps
	// dev/CI builds useful without requiring every invocation to pass ldflags.
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if Commit == "unknown" {
					Commit = s.Value
				}
			case "vcs.time":
				if BuildTime == "unknown" {
					BuildTime = s.Value
				}
			}
		}
	}

	// Some invocations (e.g. go run) do not embed VCS settings. Try git directly
	// as a last resort so local development still shows useful metadata.
	if Commit == "unknown" {
		if out, err := exec.Command("git", "rev-parse", "HEAD").Output(); err == nil {
			Commit = strings.TrimSpace(string(out))
		}
	}
	if BuildTime == "unknown" {
		if out, err := exec.Command("git", "show", "-s", "--format=%cI", "HEAD").Output(); err == nil {
			BuildTime = strings.TrimSpace(string(out))
		} else {
			BuildTime = time.Now().UTC().Format(time.RFC3339)
		}
	}
}

// String returns a human-readable version string.
func String() string {
	return fmt.Sprintf("%s (%s) built %s", Version, Commit, BuildTime)
}

// IsDev reports whether the binary was built without release metadata.
func IsDev() bool {
	return Version == "dev" || strings.HasPrefix(Version, "dev-")
}
