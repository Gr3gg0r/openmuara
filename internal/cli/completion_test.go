package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompletionCommandBash(t *testing.T) {
	cmd := newCompletionCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"bash"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "bash completion") && !strings.Contains(out, "_muara") {
		t.Errorf("expected bash completion output, got:\n%s", out)
	}
}

func TestCompletionCommandZsh(t *testing.T) {
	root := newRootCommand()
	root.AddCommand(newCompletionCommand())
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"completion", "zsh"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "#compdef muara") {
		t.Errorf("expected zsh completion output, got:\n%s", buf.String())
	}
}

func TestCompletionCommandFish(t *testing.T) {
	root := newRootCommand()
	root.AddCommand(newCompletionCommand())
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"completion", "fish"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "complete -c muara") {
		t.Errorf("expected fish completion output, got:\n%s", buf.String())
	}
}

func TestCompletionCommandPowerShell(t *testing.T) {
	cmd := newCompletionCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"powershell"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "Register-ArgumentCompleter") {
		t.Errorf("expected powershell completion output, got:\n%s", buf.String())
	}
}

func TestCompletionCommandInvalidShell(t *testing.T) {
	cmd := newCompletionCommand()
	cmd.SetArgs([]string{"tcsh"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for unsupported shell")
	}
}
