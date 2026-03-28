package checker

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
	"time"
)

// LogFunc is called for each log line during a diagnostic step.
type LogFunc func(string)

// runCmdStreaming executes a command, streaming each output line to logFn in real-time.
// The command is killed when ctx is cancelled (with a 3s grace period via WaitDelay).
// Commands are run directly (not through cmd /c) so context cancellation reliably
// kills the process.
func runCmdStreaming(ctx context.Context, logFn LogFunc, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.WaitDelay = 3 * time.Second // Force-kill if process doesn't exit after cancel
	hideConsole(cmd)                // Prevent console window flash on Windows

	r, w := io.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w

	var output strings.Builder
	done := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			output.WriteString(line + "\n")
			if logFn != nil {
				logFn(line)
			}
		}
		close(done)
	}()

	err := cmd.Run()
	w.Close()
	<-done

	return strings.TrimSpace(output.String()), err
}
