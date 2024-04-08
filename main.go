package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

// Continuous Integration tool.
// For this example, CI pipeline consists of:
// - building the program (go build) to verify the program structure is valid;
// - executing tests (go test) to ensure the program does what it's intended to do;
// - executing gofmt to ensure the program's format conforms to the standards;
// - executing git push to push the code to the remote shared Git repo.

type executer interface {
	execute() (string, error)
}

// func run takes two input params:
// proj - go project directory on which to execute the CI pipeline steps;
// out - interface to output the status of the tool.
func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}

	pipeline := make([]executer, 4)

	// builds the program to verify the structure is valid
	pipeline[0] = newExceptionStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)

	// tests with "go test -v"
	pipeline[1] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		proj,
		[]string{"test", "-v"},
	)

	// validates whether the project conforms to the Go code fotmatting standards
	// "gofmt -l" - lists files whose formatting differs from gofmt's
	pipeline[2] = newExceptionStep(
		"go fmt",
		"gofmt",
		"Gofmt: SUCCESS",
		proj,
		[]string{"-l", "."},
	)

	pipeline[3] = newTimeoutStep(
		"git push",
		"git",
		"Git Push: SUCCESS",
		proj,
		[]string{"push", "origin", "main"},
		10*time.Second,
	)

	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err //nolint:wrapcheck
		}

		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err //nolint:wrapcheck
		}
	}

	return nil
}

func main() {
	fmt.Println("Running ci pipeline...")

	proj := flag.String("p", "", "Project directory")

	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
