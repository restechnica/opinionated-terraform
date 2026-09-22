package root

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
	"github.com/restechnica/opinionated-terraform/pkg/cli/version"
	"github.com/restechnica/opinionated-terraform/pkg/core"
	"github.com/restechnica/opinionated-terraform/pkg/env"
)

var tfArgs []string

var tfCommands = []string{
	"apply", "console", "destroy", "fmt", "force-unlock",
	"graph", "import", "init", "output", "plan",
	"providers", "refresh", "show", "state", "taint",
	"test", "untaint", "validate", "workspace",
}

var otfFlags = map[string]bool{
	"--debug": true, "-d": true,
	"--verbose": true, "-v": true,
	"--help": true, "-h": true,
}

// ArgSplitter separates otf arguments from terraform passthrough arguments.
type ArgSplitter struct {
	OTFArgs []string
	TFArgs  []string
}

// Execute creates the root command and executes the CLI.
func Execute() error {
	var command = NewCommand()

	subcommands := GetSubcommandMap(command)
	splitter := NewArgSplitter(subcommands)
	tfArgs = splitter.TFArgs
	command.SetArgs(splitter.OTFArgs)

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}

	return nil
}

// NewCommand creates and returns the root command of the otf CLI.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "otf [flags] <env> <command> [terraform args...]",
		Short: "A lightweight opinionated wrapper around Terraform",
		Long: `otf wraps Terraform's partial backend configuration to make environment
switching safe and simple. It automatically re-initializes when the environment
changes and injects the right -var-file, while passing everything else straight
through to Terraform.`,
		PersistentPreRunE: persistentPreRunE,
		RunE:              runE,
		Args:              cobra.MinimumNArgs(2),
		ValidArgsFunction: completeArgs,
	}

	cmd.PersistentFlags().BoolVarP(&cli.VerboseFlag, cli.VerboseFlagName, "v", false,
		"set log level verbosity")
	cmd.PersistentFlags().BoolVarP(&cli.DebugFlag, cli.DebugFlagName, "d", false,
		"set log level verbosity to debug")

	cmd.AddCommand(version.NewCommand())

	return cmd
}

// GetSubcommandMap returns a set of all registered subcommand names and aliases for the given command,
// including cobra's internal completion commands.
func GetSubcommandMap(cmd *cobra.Command) map[string]bool {
	subcommands := map[string]bool{
		"__complete":       true,
		"__completeNoDesc": true,
	}

	for _, sub := range cmd.Commands() {
		subcommands[sub.Name()] = true
		for _, alias := range sub.Aliases {
			subcommands[alias] = true
		}
	}

	return subcommands
}

// NewArgSplitter splits os.Args using the given subcommand set.
// OTF flags and the two positional args (env, command) go to OTFArgs.
// Everything after the command position goes to TFArgs.
// If the first positional arg is a known subcommand, all remaining args go to OTFArgs.
func NewArgSplitter(subcommands map[string]bool) ArgSplitter {
	args := os.Args[1:]
	i := 0

	var otfArgs []string

	for i < len(args) {
		if otfFlags[args[i]] {
			otfArgs = append(otfArgs, args[i])
			i++
			continue
		}
		break
	}

	if i >= len(args) {
		return ArgSplitter{OTFArgs: otfArgs}
	}

	if subcommands[args[i]] {
		return ArgSplitter{OTFArgs: append(otfArgs, args[i:]...)}
	}

	remaining := args[i:]

	if len(remaining) >= 2 {
		otfArgs = append(otfArgs, remaining[0], remaining[1])
		return ArgSplitter{OTFArgs: otfArgs, TFArgs: remaining[2:]}
	}

	return ArgSplitter{OTFArgs: append(otfArgs, remaining...)}
}

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	cli.ConfigureLogging()
	return nil
}

func runE(cmd *cobra.Command, args []string) error {
	targetEnv := args[0]
	tfCommand := args[1]

	log.Debug().Str("env", targetEnv).Str("command", tfCommand).Strs("args", tfArgs).Msg("starting...")

	if err := core.Run(targetEnv, tfCommand, tfArgs); err != nil {
		return err
	}

	log.Debug().Str("env", targetEnv).Str("command", tfCommand).Msg("done")

	return nil
}

func completeArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return env.List(), cobra.ShellCompDirectiveNoFileComp
	case 1:
		return tfCommands, cobra.ShellCompDirectiveNoFileComp
	default:
		return nil, cobra.ShellCompDirectiveDefault
	}
}
