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
// If the backend config file does not exist, init runs without backend configuration.
func (api CLI) Init(env string) error {
	args := []string{"init"}

	backendConfig := filepath.Join(cli.DefaultBackendsDir, env+".tf")

	if fs.Exists(backendConfig) {
		args = append(args, "-backend-config", backendConfig, "-reconfigure")
	} else {
		log.Info().Str("path", backendConfig).Msg("backend config not found, using local state")
	}

	log.Debug().Strs("args", args).Msg("running terraform init")

	return api.Commander.Stream("terraform", args...)
}

// Run executes a terraform command with the appropriate flags for the given environment.
// All extra arguments are passed through to terraform verbatim.
func (api CLI) Run(env string, command string, extraArgs []string) error {
	args := BuildArgs(env, command, extraArgs)

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
// If the variables file does not exist on disk, -var-file is omitted.
func BuildArgs(env string, command string, extraArgs []string) []string {
	args := []string{command}

	if NeedsVarFile(command) {
		varFile := filepath.Join(cli.DefaultVariablesDir, env+".tfvars")

		if fs.Exists(varFile) {
			args = append(args, "-var-file", varFile)
		} else {
			log.Info().Str("path", varFile).Msg("variables file not found, skipping var-file injection")
		}
	}

	args = append(args, extraArgs...)

	return args
}
