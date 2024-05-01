package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
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
		"Go Build: SUCCESS\n",
		proj,
		[]string{"build", ".", "errors"},
	)

	// tests with "go test -v"
	pipeline[1] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS\n",
		proj,
		[]string{"test", "-v"},
	)

	// validates whether the project conforms to the Go code fotmatting standards
	// "gofmt -l" - lists files whose formatting differs from gofmt's
	pipeline[2] = newExceptionStep(
		"go fmt",
		"gofmt",
		"Gofmt: SUCCESS\n",
		proj,
		[]string{"-l", "."},
	)

	pipeline[3] = newTimeoutStep(
		"git push",
		"git",
		"Git Push: SUCCESS\n",
		proj,
		[]string{"push", "origin", "master"},
		10*time.Second,
	)

	// sig - a buffered channel of size one which
	// allows the app to handle at least one signal correctly
	sig := make(chan os.Signal, 1)

	errCh := make(chan error)
	done := make(chan error)

	// relay SIGINT and SIGTERM to channel sig
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// to allow concurrent execution with signal.Notify
	go func() {
		for _, s := range pipeline {
			msg, err := s.execute()
			if err != nil {
				errCh <- err
				return
			}

			_, err = fmt.Fprintln(out, msg)
			if err != nil {
				errCh <- err
				return
			}
		}

		close(done)
	}()

	for {
		select {
		case rec := <-sig:
			signal.Stop(sig)
			return fmt.Errorf("%s: exiting: %w", rec, ErrSignal)
		case err := <-errCh:
			return err
		case <-done:
			return nil
		}
	}
}

func main() {
	fmt.Println("Running ci pipeline...")

	proj := flag.String("p", "", "Project directory")

	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
