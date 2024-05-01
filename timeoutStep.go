package main

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

// timeoutStep extends step type by adding context.
type timeoutStep struct {
	step
	timeout time.Duration
}

func newTimeoutStep(name, exe, message, proj string,
	args []string, timeout time.Duration) timeoutStep {

	s := timeoutStep{}

	s.step = newStep(name, exe, message, proj, args)
	s.timeout = timeout

	if s.timeout == 0 {
		s.timeout = 30 * time.Second
	}

	return s
}

func (s timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, s.exe, s.args...) //nolint:gosec
	cmd.Dir = s.proj

	var out bytes.Buffer
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", &stepError{
				step:  s.name,
				msg:   "failed time out: " + out.String(),
				cause: context.DeadlineExceeded,
			}
		}

		cmd.Stdout = &out

		return "", &stepError{
			step:  s.name,
			msg:   "failed to execute: " + out.String(),
			cause: err,
		}
	}

	return s.message, nil
}
