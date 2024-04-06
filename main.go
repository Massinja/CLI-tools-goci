package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Continuous Integration tool.
// For this example, CI pipeline consists of:
// - building the program (go build) to verify the program structure is valid;
// - executing tests (go test) to ensure the program does what it's intended to do;
// - executing gofmt to ensure the program's format conforms to the standards;
// - executing git push to push the code to the remote shared Git repo.

// func run takes two input params:
// proj - go project directory on which to execute the CI pipeline steps;
// out - interface to output the status of the tool.
func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}

	pipeline := make([]step, 1)
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)
	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}

	return nil //nolint:wrapcheck
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
