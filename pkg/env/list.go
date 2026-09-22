package env

import (
	"path/filepath"
	"strings"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
)

// List returns the available environment names by listing backend config files.
func List() []string {
	matches, err := filepath.Glob(filepath.Join(cli.DefaultBackendsDir, "*.tf"))
	if err != nil {
		return nil
	}

	envs := make([]string, 0, len(matches))
	for _, match := range matches {
		name := strings.TrimSuffix(filepath.Base(match), ".tf")
		envs = append(envs, name)
	}

	return envs
}
