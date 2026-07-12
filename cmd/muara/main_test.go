package main

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestMainVersion(_ *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"muara", "version"}

	main()
}

func TestMainExecuteError(t *testing.T) {
	oldExecute := execute
	oldExit := exitFunc
	defer func() {
		execute = oldExecute
		exitFunc = oldExit
	}()

	execute = func(context.Context) error { return errors.New("boom") }
	var exitCode int
	exitFunc = func(code int) { exitCode = code }

	main()

	if exitCode != 1 {
		t.Errorf("exit code: want 1, got %d", exitCode)
	}
}
