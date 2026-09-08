package core

import (
	"github.com/rs/zerolog/log"

	"github.com/restechnica/opinionated-terraform/pkg/env"
	"github.com/restechnica/opinionated-terraform/pkg/terraform"
)

// Run validates the environment, initializes if needed, and executes a terraform command.
func Run(tf terraform.API, target string, command string, extraArgs []string) (err error) {
	if err = env.Validate(target); err != nil {
		return err
	}

	var current string

	if current, err = env.ReadCurrent(); err != nil {
		return err
	}

	if current != target {
		log.Info().Str("from", current).Str("to", target).Msg("environment changed, running init")

		if err = tf.Init(target); err != nil {
			return err
		}

		if err = env.WriteCurrent(target); err != nil {
			return err
		}
	} else {
		log.Debug().Str("env", target).Msg("environment unchanged, skipping init")
	}

	return tf.Run(target, command, extraArgs)
}
