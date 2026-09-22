package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
)

func TestReadCurrent(t *testing.T) {
	t.Run("ReturnEmptyStringsWhenFileDoesNotExist", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		env, hash, err := ReadCurrent()

		assert.NoError(t, err)
		assert.Empty(t, env)
		assert.Empty(t, hash)
	})

	t.Run("ReturnEnvAndHashWhenFileExists", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(".terraform", 0755))
		require.NoError(t, os.WriteFile(filepath.Join(".terraform", cli.DefaultEnvFile), []byte("staging\nabc123\n"), 0644))

		env, hash, err := ReadCurrent()

		assert.NoError(t, err)
		assert.Equal(t, "staging", env)
		assert.Equal(t, "abc123", hash)
	})

	t.Run("ReturnEmptyHashForOldFormat", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(".terraform", 0755))
		require.NoError(t, os.WriteFile(filepath.Join(".terraform", cli.DefaultEnvFile), []byte("prod\n"), 0644))

		env, hash, err := ReadCurrent()

		assert.NoError(t, err)
		assert.Equal(t, "prod", env)
		assert.Empty(t, hash)
	})
}

func TestWriteCurrent(t *testing.T) {
	t.Run("WriteEnvAndHashToFile", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(".terraform", 0755))

		err := WriteCurrent("prod", "abc123")
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(".terraform", cli.DefaultEnvFile))
		require.NoError(t, err)

		assert.Equal(t, "prod\nabc123\n", string(data))
	})
}
