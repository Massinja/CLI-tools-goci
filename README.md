Continuous Integration tool.  
For this example, CI pipeline consists of:  
- building the program (go build) to verify the program structure is valid;  
- executing tests (go test) to ensure the program does what it's intended to do;  
- executing gofmt to ensure the program's format conforms to the standards;  
- executing git push to push the code to the remote shared Git repo.  
