package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func fzf(root string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fzfCmd := exec.CommandContext(ctx, "fzf")
	findCmd := exec.CommandContext(ctx, "find", root, "-type", "d")

	// show errors if fails
	findCmd.Stderr = os.Stderr
	fzfCmd.Stderr = os.Stderr

	// Piping
	r, w := io.Pipe()
	findCmd.Stdout = w
	fzfCmd.Stdin = r
	defer r.Close()
	defer w.Close()

	// Fzf stdout
	var res strings.Builder
	fzfCmd.Stdout = &res

	// Can also use goroutine and Run, instead of Start and Wait for async use
	if err := findCmd.Start(); err != nil {
		return "", fmt.Errorf("find command failed to start: %w", err)
	}
	if err := fzfCmd.Start(); err != nil {
		return "", fmt.Errorf("fzf command failed to start: %w", err)
	}

	findCh := make(chan error, 1) // ignored
	go func() {
		findCh <- findCmd.Wait()
		w.Close()
	}()

	fzfCh := make(chan error, 1)
	go func() {
		fzfCh <- fzfCmd.Wait()
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-fzfCh:
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
				return "", nil // gracefully cancelled
			}
			return "", fmt.Errorf("fzf command failed: %w", err)
		}
	}
	return strings.TrimSpace(res.String()), nil
}
