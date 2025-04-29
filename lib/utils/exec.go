package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ExecCommand runs the given command in a process,
// captures the stdout and stderr in a combined output
// and returns an error message, that's printing the full
// contents of stderr on top of the error message itself.
func ExecCommand(name string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	command := exec.Command(name, args...)
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return stdout.String(), fmt.Errorf("%w; caputed stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
