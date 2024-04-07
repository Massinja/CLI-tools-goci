package main

import "os/exec"

// A pipeline step.
type step struct {
	name    string   // step name
	exe     string   // executable name of the external tool we want to execute
	args    []string // arguments for the executable
	message string   // output message in case of success
	proj    string   // target project on which to execute the task
}

func (s step) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &stepError{
			step:  s.name,
			msg:   "failed to execute\n",
			cause: err,
		}
	}

	return s.message, nil
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
