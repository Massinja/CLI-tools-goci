package main

import (
	"bytes"
	"os/exec"
)

// A pipeline step.
type step struct {
	name    string   // step name
	exe     string   // executable name of the external tool we want to execute
	args    []string // arguments for the executable
	message string   // output message in case of success
	proj    string   // target project on which to execute the task
}

// newStep instantiates and returns new step.
func newStep(name, exe, message, proj string, args []string) step {
	return step{
		name:    name,
		exe:     exe,
		args:    args,
		message: message,
		proj:    proj,
	}
}

func (s step) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...) //nolint:gosec
	cmd.Dir = s.proj

	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", &stepError{
			step:  s.name,
			msg:   "failed to execute: " + stderr.String(),
			cause: err,
		}
	}

	if s.name == "go fmt" {
		// gofmt -l will list unformatted files
		if stdout.Len() > 0 {
			return "", &stepError{
				step: s.name,
				msg:  "invalid format: " + stdout.String(),
			}
		}
	}

	return s.message, nil
}
