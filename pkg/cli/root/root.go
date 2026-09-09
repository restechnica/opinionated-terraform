package root

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
	"github.com/restechnica/opinionated-terraform/pkg/cli/version"
	"github.com/restechnica/opinionated-terraform/pkg/core"
	"github.com/restechnica/opinionated-terraform/pkg/terraform"
)

// terraformExtraArgs holds the raw terraform arguments extracted before cobra parses.
var terraformExtraArgs []string

// otfFlags is the set of flags that belong to otf and should be parsed by cobra.
var otfFlags = map[string]bool{
	"--debug": true, "-d": true,
	"--verbose": true, "-v": true,
	"--help": true, "-h": true,
}

// Execute creates the root command and executes the CLI.
func Execute() error {
	var command = NewCommand()

	subcommands := make(map[string]bool)
	for _, sub := range command.Commands() {
		subcommands[sub.Name()] = true
		for _, alias := range sub.Aliases {
			subcommands[alias] = true
		}
	}

	otfArgs, tfArgs := splitArgs(os.Args[1:], subcommands)
	terraformExtraArgs = tfArgs
	command.SetArgs(otfArgs)

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}

	return nil
}

// splitArgs separates otf arguments from terraform passthrough arguments.
// OTF flags and the two positional args (env, command) go to cobra.
// Everything after the command position is passed to terraform verbatim.
// If the first positional arg is a known subcommand, all remaining args go to cobra.
func splitArgs(args []string, subcommands map[string]bool) (otfArgs []string, tfArgs []string) {
	i := 0

	for i < len(args) {
		if otfFlags[args[i]] {
			otfArgs = append(otfArgs, args[i])
			i++
			continue
		}
		break
	}

	if i >= len(args) {
		return otfArgs, nil
	}

	if subcommands[args[i]] {
		return append(otfArgs, args[i:]...), nil
	}

	remaining := args[i:]

	if len(remaining) >= 2 {
		otfArgs = append(otfArgs, remaining[0], remaining[1])
		tfArgs = remaining[2:]
	} else {
		otfArgs = append(otfArgs, remaining...)
	}

	return otfArgs, tfArgs
}

// NewCommand creates and returns the root command of the otf CLI.
func NewCommand() *cobra.Command {
	tf := terraform.NewCLI()

	cmd := &cobra.Command{
		Use:   "otf [flags] <env> <command> [terraform args...]",
		Short: "A lightweight opinionated wrapper around Terraform",
		Long: `otf wraps Terraform's partial backend configuration to make environment
switching safe and simple. It automatically re-initializes when the environment
changes and injects the right -var-file, while passing everything else straight
through to Terraform.`,
		PersistentPreRunE: persistentPreRunE,
		RunE:              newRunE(tf),
		Args:              cobra.MinimumNArgs(2),
	}

	cmd.PersistentFlags().BoolVarP(&cli.VerboseFlag, cli.VerboseFlagName, "v", false,
		"set log level verbosity")
	cmd.PersistentFlags().BoolVarP(&cli.DebugFlag, cli.DebugFlagName, "d", false,
		"set log level verbosity to debug")

	cmd.AddCommand(version.NewCommand())

	return cmd
}

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	cli.ConfigureLogging()
	return nil
}

func newRunE(tf terraform.API) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		env := args[0]
		tfCommand := args[1]
		tfArgs := terraformExtraArgs

		log.Debug().Str("env", env).Str("command", tfCommand).Strs("args", tfArgs).Msg("starting...")

		if err := core.Run(tf, env, tfCommand, tfArgs); err != nil {
			return err
		}

		log.Debug().Str("env", env).Str("command", tfCommand).Msg("done")

		return nil
	}
}
