package cli

import (
	"context"
	"os"
	"testing"
)

func TestExecuteVersion(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"muara", "version"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := Execute(ctx); err != nil {
		t.Fatalf("execute version: %v", err)
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"muara", "not-a-command"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := Execute(ctx); err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestExecuteVersionFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"muara", "--version"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := Execute(ctx); err != nil {
		t.Fatalf("execute --version: %v", err)
	}
}

func TestExecuteVersionShortFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"muara", "-v"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := Execute(ctx); err != nil {
		t.Fatalf("execute -v: %v", err)
	}
}
