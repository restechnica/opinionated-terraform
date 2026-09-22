package env

import (
	"path/filepath"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
)

// GetBackendFilePath returns the backend config file path for the given environment.
func GetBackendFilePath(env string) string {
	return filepath.Join(cli.DefaultBackendsDir, env+".tf")
}

// GetVariablesFilePath returns the variables file path for the given environment.
func GetVariablesFilePath(env string) string {
	return filepath.Join(cli.DefaultVariablesDir, env+".tfvars")
}
