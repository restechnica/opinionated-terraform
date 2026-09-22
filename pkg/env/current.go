package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
)

// ReadCurrent reads the current environment and backend config hash from the tracking file.
// Returns empty strings if the file does not exist or uses the old single-line format.
func ReadCurrent() (env string, backendHash string, err error) {
	envFile := filepath.Join(".terraform", cli.DefaultEnvFile)

	bytes, err := os.ReadFile(envFile)
	if os.IsNotExist(err) {
		return "", "", nil
	}

	if err != nil {
		return "", "", fmt.Errorf("reading env file: %w", err)
	}

	lines := strings.SplitN(strings.TrimSpace(string(bytes)), "\n", 2)

	env = strings.TrimSpace(lines[0])

	if len(lines) > 1 {
		backendHash = strings.TrimSpace(lines[1])
	}

	return env, backendHash, nil
}

// WriteCurrent writes the environment name and backend config hash to the tracking file.
func WriteCurrent(env string, backendHash string) error {
	envFile := filepath.Join(".terraform", cli.DefaultEnvFile)

	content := env + "\n" + backendHash + "\n"

	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	return nil
}
