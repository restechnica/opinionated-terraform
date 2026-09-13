package terraform

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
	"github.com/restechnica/opinionated-terraform/pkg/commander"
	"github.com/restechnica/opinionated-terraform/pkg/fs"
)

// varFileCommands is the set of terraform commands that accept -var-file.
var varFileCommands = map[string]bool{
	"plan":    true,
	"apply":   true,
	"destroy": true,
	"refresh": true,
	"import":  true,
	"console": true,
}

// CLI is a terraform.API to interact with the Terraform CLI.
type CLI struct {
	Commander commander.Commander
}

// NewCLI creates a new CLI with a commander to run terraform commands.
// Returns the new CLI.
func NewCLI() *CLI {
	return &CLI{Commander: commander.NewExecCommander()}
}

// Init runs terraform init with the backend config for the given environment.
// If the backends directory does not exist, init runs without backend configuration.
// If the backends directory exists but the env-specific file is missing, an error is returned.
func (api CLI) Init(env string) error {
	args := []string{"init"}

	backendConfig := filepath.Join(cli.DefaultBackendsDir, env+".tf")

	if fs.Exists(cli.DefaultBackendsDir) {
		if !fs.Exists(backendConfig) {
			return fmt.Errorf("backend config %q not found", backendConfig)
		}

		args = append(args, "-backend-config", backendConfig, "-reconfigure")
	} else {
		log.Info().Msg("no backends directory found, using local state")
	}

	log.Debug().Strs("args", args).Msg("running terraform init")

	return api.Commander.Stream("terraform", args...)
}

// Run executes a terraform command with the appropriate flags for the given environment.
// All extra arguments are passed through to terraform verbatim.
func (api CLI) Run(env string, command string, extraArgs []string) error {
	args, err := BuildArgs(env, command, extraArgs)
	if err != nil {
		return err
	}

	log.Debug().Strs("args", args).Msg("running terraform")

	if err := api.Commander.Stream("terraform", args...); err != nil {
		var exitErr *exec.ExitError

		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		return fmt.Errorf("running terraform: %w", err)
	}

	return nil
}

// NeedsVarFile returns true if the terraform command accepts -var-file.
func NeedsVarFile(command string) bool {
	return varFileCommands[command]
}

// BuildArgs constructs the full terraform argument slice, injecting -var-file when appropriate.
// If the variables directory does not exist, -var-file is omitted.
// If the variables directory exists but the env-specific file is missing, an error is returned.
func BuildArgs(env string, command string, extraArgs []string) ([]string, error) {
	args := []string{command}

	if NeedsVarFile(command) {
		varFile := filepath.Join(cli.DefaultVariablesDir, env+".tfvars")

		if fs.Exists(cli.DefaultVariablesDir) {
			if !fs.Exists(varFile) {
				return nil, fmt.Errorf("variables file %q not found", varFile)
			}

			args = append(args, "-var-file", varFile)
		} else {
			log.Info().Msg("no variables directory found, skipping var-file injection")
		}
	}

	args = append(args, extraArgs...)

	return args, nil
}
