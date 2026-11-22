// Execute a shell command in the project workspace.
package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/WywySenarios/wc-tools/shared"
)

const (
	defaultTimeout = 60 * time.Second
	maxTimeout     = 24 * time.Hour

	exitUsage  = 2
	exitNotRun = 127
	exitTimed  = 124
	exitCancel = 130
)

type input struct {
	Command string   `json:"command"`
	Timeout *float64 `json:"timeout"`
	Workdir string   `json:"workdir"`
}

func main() {
	os.Exit(run())
}

func run() int {
	var in input
	if err := shared.ReadArgs(os.Args[1:], os.Stdin, &in); err != nil {
		fmt.Fprintf(os.Stderr, "bash: %v\n", err)
		return exitUsage
	}

	if strings.TrimSpace(in.Command) == "" {
		fmt.Fprintln(os.Stderr, "bash: command is required")
		return exitUsage
	}

	timeout, err := resolveTimeout(in.Timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bash: %v\n", err)
		return exitUsage
	}

	workdir, err := resolveWorkdir(in.Workdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bash: %v\n", err)
		return exitUsage
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", in.Command)
	cmd.Dir = workdir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// Kill the whole process group so background children cannot hold the output pipes open.
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	cmd.WaitDelay = time.Second

	err = cmd.Run()
	switch {
	case err == nil:
		return 0
	case ctx.Err() == context.DeadlineExceeded:
		fmt.Fprintf(os.Stderr, "bash: timed out after %s\n", timeout)
		return exitTimed
	case ctx.Err() == context.Canceled:
		return exitCancel
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			return 128 + int(status.Signal())
		}
		if code := exitErr.ExitCode(); code > 0 {
			return code
		}
		return 1
	}

	fmt.Fprintf(os.Stderr, "bash: %v\n", err)
	return exitNotRun
}

func resolveTimeout(seconds *float64) (time.Duration, error) {
	if seconds == nil {
		return defaultTimeout, nil
	}
	if math.IsNaN(*seconds) || math.IsInf(*seconds, 0) {
		return 0, errors.New("timeout must be a finite number of seconds")
	}
	if *seconds < 1 {
		return 0, errors.New("timeout must be at least 1 second")
	}
	if *seconds > maxTimeout.Seconds() {
		return 0, fmt.Errorf("timeout must not exceed %v", maxTimeout)
	}
	return time.Duration(*seconds * float64(time.Second)), nil
}

func resolveWorkdir(dir string) (string, error) {
	if dir == "" {
		return "", nil
	}
	cleaned := filepath.Clean(dir)
	if filepath.IsAbs(cleaned) {
		return "", errors.New("workdir must be relative to the project root")
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", errors.New("workdir escapes the project root")
	}
	return cleaned, nil
}
