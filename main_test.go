package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		proj     string
		out      string
		expError error
	}{
		{
			name:     "success",
			proj:     "./testdata/tool/",
			out:      "Go Build: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCCESS\nGit Push: SUCCESS\n",
			expError: nil,
		},
		{
			name:     "fail",
			proj:     "./testdata/toolError/",
			out:      "",
			expError: &stepError{step: "go build"},
		},
		{
			name:     "failformat",
			proj:     "./testdata/toolFmtErr/",
			out:      "",
			expError: &stepError{step: "go fmt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := run(tc.proj, &out)

			if tc.expError != nil {
				if err == nil {
					t.Errorf("expected error: %v; got nil instead", tc.expError)
					return
				}

				if !errors.Is(tc.expError, err) {
					t.Errorf("expected error: %v; got %v instead", tc.expError, err)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if out.String() != tc.out {
				t.Errorf("expected output: %s; got: %s", tc.out, out.String())
			}
		})
	}
}
