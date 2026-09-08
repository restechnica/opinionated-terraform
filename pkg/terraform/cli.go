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
func (api CLI) Init(env string) error {
	backendConfig := filepath.Join(cli.DefaultBackendsDir, env+".tf")

	args := []string{"init", "-backend-config", backendConfig, "-reconfigure"}

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
func BuildArgs(env string, command string, extraArgs []string) []string {
	args := []string{command}

	if NeedsVarFile(command) {
		varFile := filepath.Join(cli.DefaultVariablesDir, env+".tfvars")
		args = append(args, "-var-file", varFile)
	}

	args = append(args, extraArgs...)

	return args
}
