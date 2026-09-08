package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
)

func TestValidate(t *testing.T) {
	t.Run("SucceedWhenBothFilesExist", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(cli.DefaultBackendsDir, 0755))
		require.NoError(t, os.MkdirAll(cli.DefaultVariablesDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultBackendsDir, "prod.tf"), []byte(""), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultVariablesDir, "prod.tfvars"), []byte(""), 0644))

		var got = Validate("prod")

		assert.NoError(t, got)
	})

	t.Run("FailWhenBackendConfigIsMissing", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(cli.DefaultVariablesDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultVariablesDir, "prod.tfvars"), []byte(""), 0644))

		var got = Validate("prod")

		assert.ErrorContains(t, got, "backend config not found")
	})

	t.Run("FailWhenVariablesFileIsMissing", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(cli.DefaultBackendsDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultBackendsDir, "prod.tf"), []byte(""), 0644))

		var got = Validate("prod")

		assert.ErrorContains(t, got, "variables file not found")
	})
}

func TestReadCurrent(t *testing.T) {
	t.Run("ReturnEmptyStringWhenFileDoesNotExist", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		var got, err = ReadCurrent()

		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("ReturnEnvNameWhenFileExists", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(".terraform", 0755))
		require.NoError(t, os.WriteFile(filepath.Join(".terraform", cli.DefaultEnvFile), []byte("staging\n"), 0644))

		var want = "staging"
		var got, err = ReadCurrent()

		assert.NoError(t, err)
		assert.Equal(t, want, got, `want: '%s', got: '%s'`, want, got)
	})
}

func TestWriteCurrent(t *testing.T) {
	t.Run("WriteEnvToFile", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(".terraform", 0755))

		var err = WriteCurrent("prod")
		require.NoError(t, err)

		var want = "prod\n"
		data, err := os.ReadFile(filepath.Join(".terraform", cli.DefaultEnvFile))
		require.NoError(t, err)

		var got = string(data)

		assert.Equal(t, want, got, `want: '%s', got: '%s'`, want, got)
	})
}
