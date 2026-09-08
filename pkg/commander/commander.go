package commander

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// Commander is an interface for executing CLI commands.
// It extends the go-cmder pattern with a Stream method for interactive commands
// that need stdin/stdout/stderr passthrough.
type Commander interface {
	// Output runs a command and returns its combined output as a string.
	Output(name string, arg ...string) (string, error)

	// Run runs a command and returns only an error.
	Run(name string, arg ...string) error

	// Stream runs a command with stdin/stdout/stderr connected to the terminal.
	// Use this for interactive commands that need user input (e.g. terraform apply).
	Stream(name string, arg ...string) error
}

// ExecCommander is the production implementation using os/exec.
type ExecCommander struct{}

// NewExecCommander creates a new ExecCommander.
func NewExecCommander() *ExecCommander {
	return &ExecCommander{}
}

// Output runs a command and returns its combined stdout and stderr.
func (c ExecCommander) Output(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		return buf.String(), fmt.Errorf("command [%s %v] failed: %w", name, arg, err)
	}

	return buf.String(), nil
}

// Run runs a command and discards its output.
func (c ExecCommander) Run(name string, arg ...string) error {
	_, err := c.Output(name, arg...)
	return err
}

// Stream runs a command with stdin/stdout/stderr connected to the terminal.
func (c ExecCommander) Stream(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
