package core

import (
	"github.com/rs/zerolog/log"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
	"github.com/restechnica/opinionated-terraform/pkg/env"
	"github.com/restechnica/opinionated-terraform/pkg/fs"
	"github.com/restechnica/opinionated-terraform/pkg/terraform"
)

// Run initializes if needed and executes a terraform command for the given environment.
func Run(target string, command string, tfArgs []string) (err error) {
	var currentEnv, currentHash, targetHash string

	if currentEnv, currentHash, err = env.ReadCurrent(); err != nil {
		return err
	}

	if targetHash, err = backendHash(target); err != nil {
		return err
	}

	var tf = terraform.NewCLI()

	envHasChanged := currentEnv != target
	backendIsModified := currentHash != targetHash

	if !envHasChanged && !backendIsModified {
		log.Debug().Str("env", target).Msg("environment unchanged, skipping init")
		return tf.Run(target, command, tfArgs)
	}

	if envHasChanged {
		log.Info().Str("from", currentEnv).Str("to", target).Msg("environment changed, running init")
	} else {
		log.Info().Str("env", target).Msg("backend config changed, running init")
	}

	if err = tf.Init(target); err != nil {
		return err
	}

	if err = env.WriteCurrent(target, targetHash); err != nil {
		return err
	}

	return tf.Run(target, command, tfArgs)
}

// backendHash returns the SHA-256 hash of the backend config file for the given environment.
// Returns an empty string if the backends directory does not exist.
func backendHash(target string) (string, error) {
	if !fs.Exists(cli.DefaultBackendsDir) {
		return "", nil
	}

	return fs.HashFile(env.GetBackendFilePath(target))
}
