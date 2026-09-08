package root

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
	"github.com/restechnica/opinionated-terraform/pkg/cli/version"
	"github.com/restechnica/opinionated-terraform/pkg/commander"
	"github.com/restechnica/opinionated-terraform/pkg/core"
)

// Execute creates the root command and executes the CLI
func Execute() error {
	var command = NewCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}

	return nil
}

// NewCommand creates and returns the root command of the otf CLI.
func NewCommand() *cobra.Command {
	cmdr := commander.NewExecCommander()

	cmd := &cobra.Command{
		Use:   "otf <env> <command> [terraform args...]",
		Short: "A lightweight opinionated wrapper around Terraform",
		Long: `otf wraps Terraform's partial backend configuration to make environment
switching safe and simple. It automatically re-initializes when the environment
changes and injects the right -var-file, while passing everything else straight
through to Terraform.`,
		PersistentPreRunE:  persistentPreRunE,
		RunE:               newRunE(cmdr),
		Args:               cobra.MinimumNArgs(2),
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
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

func newRunE(cmdr commander.Commander) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		env := args[0]
		tfCommand := args[1]
		tfArgs := args[2:]

		log.Debug().Str("env", env).Str("command", tfCommand).Strs("args", tfArgs).Msg("starting...")

		if err := core.ValidateEnv(env); err != nil {
			return err
		}

		if err := core.InitIfNeeded(cmdr, env); err != nil {
			return err
		}

		if err := core.Run(cmdr, env, tfCommand, tfArgs); err != nil {
			return err
		}

		log.Debug().Str("env", env).Str("command", tfCommand).Msg("done")

		return nil
	}
}
