package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
)

// Validate checks that the backend config and variable files exist for the given environment.
func Validate(env string) error {
	backendFile := filepath.Join(cli.DefaultBackendsDir, env+".tf")
	variablesFile := filepath.Join(cli.DefaultVariablesDir, env+".tfvars")

	if _, err := os.Stat(backendFile); os.IsNotExist(err) {
		return fmt.Errorf("backend config not found: %s", backendFile)
	}

	if _, err := os.Stat(variablesFile); os.IsNotExist(err) {
		return fmt.Errorf("variables file not found: %s", variablesFile)
	}

	return nil
}

// ReadCurrent reads the current environment from the .terraform/.otf tracking file.
// Returns an empty string if the file does not exist.
func ReadCurrent() (data string, err error) {
	var bytes []byte

	envFile := filepath.Join(".terraform", cli.DefaultEnvFile)

	if bytes, err = os.ReadFile(envFile); os.IsNotExist(err) {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("reading env file: %w", err)
	}

	return strings.TrimSpace(string(bytes)), nil
}

// WriteCurrent writes the environment name to the .terraform/.otf tracking file.
func WriteCurrent(env string) error {
	envFile := filepath.Join(".terraform", cli.DefaultEnvFile)

	if err := os.WriteFile(envFile, []byte(env+"\n"), 0644); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	return nil
}
