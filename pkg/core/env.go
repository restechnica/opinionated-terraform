package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
	"github.com/restechnica/opinionated-terraform/pkg/commander"
)

// ValidateEnv checks that the backend config and variable files exist for the given environment.
func ValidateEnv(env string) error {
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

// ReadCurrentEnv reads the current environment from the .terraform/.otf tracking file.
// Returns an empty string if the file does not exist.
func ReadCurrentEnv() (data string, err error) {
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

// WriteCurrentEnv writes the environment name to the .terraform/.otf tracking file.
func WriteCurrentEnv(env string) error {
	envFile := filepath.Join(".terraform", cli.DefaultEnvFile)

	if err := os.WriteFile(envFile, []byte(env+"\n"), 0644); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	return nil
}

// InitIfNeeded runs terraform init with the appropriate backend config if the environment
// has changed since the last init. It always inits if no previous environment is tracked.
func InitIfNeeded(cmdr commander.Commander, env string) error {
	current, err := ReadCurrentEnv()
	if err != nil {
		return err
	}

	if current == env {
		log.Debug().Str("env", env).Msg("environment unchanged, skipping init")
		return nil
	}

	log.Info().Str("from", current).Str("to", env).Msg("environment changed, running init")

	if err := runInit(cmdr, env); err != nil {
		return err
	}

	return WriteCurrentEnv(env)
}

// runInit executes terraform init with the backend config for the given environment.
func runInit(cmdr commander.Commander, env string) error {
	backendConfig := filepath.Join(cli.DefaultBackendsDir, env+".tf")

	args := []string{"init", "-backend-config", backendConfig, "-reconfigure"}

	log.Debug().Strs("args", args).Msg("running terraform init")

	return cmdr.Stream("terraform", args...)
}
