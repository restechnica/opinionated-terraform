package core

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

// Run executes a terraform command with the appropriate flags for the given environment.
// All extra arguments are passed through to terraform verbatim.
func Run(cmdr commander.Commander, env string, command string, extraArgs []string) error {
	args := BuildArgs(env, command, extraArgs)

	log.Debug().Strs("args", args).Msg("running terraform")

	if err := cmdr.Stream("terraform", args...); err != nil {
		var exitErr *exec.ExitError

		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		return fmt.Errorf("running terraform: %w", err)
	}

	return nil
}
