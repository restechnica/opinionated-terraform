package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
)

func TestValidateEnv(t *testing.T) {
	t.Run("succeeds when both files exist", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(cli.DefaultBackendsDir, 0755))
		require.NoError(t, os.MkdirAll(cli.DefaultVariablesDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultBackendsDir, "prod.tf"), []byte(""), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultVariablesDir, "prod.tfvars"), []byte(""), 0644))

		err := ValidateEnv("prod")
		assert.NoError(t, err)
	})

	t.Run("fails when backend config is missing", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(cli.DefaultVariablesDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultVariablesDir, "prod.tfvars"), []byte(""), 0644))

		err := ValidateEnv("prod")
		assert.ErrorContains(t, err, "backend config not found")
	})

	t.Run("fails when variables file is missing", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(cli.DefaultBackendsDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultBackendsDir, "prod.tf"), []byte(""), 0644))

		err := ValidateEnv("prod")
		assert.ErrorContains(t, err, "variables file not found")
	})
}

func TestReadCurrentEnv(t *testing.T) {
	t.Run("returns empty string when file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))

		env, err := ReadCurrentEnv()
		assert.NoError(t, err)
		assert.Empty(t, env)
	})

	t.Run("returns env name when file exists", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(".terraform", 0755))
		require.NoError(t, os.WriteFile(filepath.Join(".terraform", cli.DefaultEnvFile), []byte("staging\n"), 0644))

		env, err := ReadCurrentEnv()
		assert.NoError(t, err)
		assert.Equal(t, "staging", env)
	})
}

func TestWriteCurrentEnv(t *testing.T) {
	t.Run("writes env to file", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(".terraform", 0755))

		err := WriteCurrentEnv("prod")
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(".terraform", cli.DefaultEnvFile))
		require.NoError(t, err)
		assert.Equal(t, "prod\n", string(data))
	})
}
