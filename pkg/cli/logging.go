package cli

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ConfigureLogging sets up zerolog with console output and the appropriate log level
func ConfigureLogging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	SetLogLevel()
}

// SetLogLevel sets the global log level based on flags
func SetLogLevel() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	if VerboseFlag {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if DebugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
