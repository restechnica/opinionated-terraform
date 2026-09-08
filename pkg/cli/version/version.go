package version

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/restechnica/opinionated-terraform/pkg/core"
)

// NewCommand creates and returns the version command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE:  runE,
	}
}

func runE(cmd *cobra.Command, args []string) error {
	log.Debug().Str("command", "version").Msg("starting...")

	if err := core.Version(); err != nil {
		return err
	}

	log.Debug().Str("command", "version").Msg("done")

	return nil
}
