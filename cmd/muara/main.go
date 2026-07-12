// Command muara is the CLI entry point for the OpenMuara payment virtualization layer.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gr3gg0r/openmuara/internal/cli"
)

// execute and exitFunc are variables so tests can intercept them without
// spawning a subprocess.
var (
	execute  = cli.Execute
	exitFunc = os.Exit
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := execute(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitFunc(1)
	}
}
