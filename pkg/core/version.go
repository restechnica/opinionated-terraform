package core

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/rs/zerolog/log"

	"github.com/restechnica/opinionated-terraform/internal/ldflags"
)

func Version() error {
	log.Debug().Str("command", "version").Msg("starting...")

	var info *debug.BuildInfo
	var ok bool

	if info, ok = debug.ReadBuildInfo(); !ok {
		return errors.New("failed to read build info")
	}

	var arch, os string

	for _, setting := range info.Settings {
		switch setting.Key {
		case "GOARCH":
			arch = setting.Value
		case "GOOS":
			os = setting.Value
		}
	}

	fmt.Printf("otf %s %s %s/%s\n", ldflags.Version, info.GoVersion, os, arch)

	log.Debug().Str("command", "version").Msg("done")

	return nil
}
