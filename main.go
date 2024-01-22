package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

/*
Continous Integration tool.
For this example, CI pipeline consists of:
- building the program (go build) to verify the program structure is valid
- executing tests (go test) to ensure the program does what it's intended to do
- executing gofmt to ensure the program's format conforms to the standards
- executing git push to push the code to the remote shared Git repo
*/

// func run takes two input params:
// proj - go project directory on which to execute the CI pipeline steps
// out - interface to output the status of the tool
func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("Project directory is missing")
	}

	// go build doesn't create an exec file
	// when building multiple packages at the same time
	// use any package from Go's standlib (i.e. "errors")
	// to avoid clean up
	args := []string{"build", ".", "errors"}
	cmd := exec.Command("go", args...)
	cmd.Dir = proj
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("'go build' failed: %s", err)
	}
	_, err := fmt.Fprintln(out, "Go Build: SUCCESS")
	return err

}

func main(){
	proj := flag.String("p", "", "Project directory")
	flag.Parse()
	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
