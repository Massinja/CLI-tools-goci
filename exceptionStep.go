package main

import (
	"bytes"
	"os/exec"
)

type exceptionStep struct {
	step
}

func newExceptionStep(name, exe, message, proj string, args []string) exceptionStep {
	s := exceptionStep{}
	s.step = newStep(name, exe, message, proj, args)

	return s
}

// extends (step).execute method to be able to handle program output.
func (s exceptionStep) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &stepError{
			step:  s.name,
			msg:   "failed to execute:\n" + out.String(),
			cause: err,
		}
	}

	if out.Len() > 0 {
		return "", &stepError{
			step: s.name,
			msg:  "invalid format: " + out.String(),
		}
	}

	return s.message, nil
}
