package env

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
)

func TestGetBackendFilePath(t *testing.T) {
	var got = GetBackendFilePath("prod")
	var want = filepath.Join(cli.DefaultBackendsDir, "prod.tf")

	assert.Equal(t, want, got)
}

func TestGetVariablesFilePath(t *testing.T) {
	var got = GetVariablesFilePath("prod")
	var want = filepath.Join(cli.DefaultVariablesDir, "prod.tfvars")

	assert.Equal(t, want, got)
}
