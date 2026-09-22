package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
)

func TestList(t *testing.T) {
	t.Run("ReturnEmptyWhenNoBackendsDir", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))

		got := List()

		assert.Empty(t, got)
	})

	t.Run("ReturnEnvironmentNames", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(cli.DefaultBackendsDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultBackendsDir, "dev.tf"), []byte{}, 0644))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultBackendsDir, "staging.tf"), []byte{}, 0644))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultBackendsDir, "prod.tf"), []byte{}, 0644))

		got := List()

		assert.Len(t, got, 3)
		assert.Contains(t, got, "dev")
		assert.Contains(t, got, "staging")
		assert.Contains(t, got, "prod")
	})

	t.Run("IgnoreNonTfFiles", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))

		require.NoError(t, os.MkdirAll(cli.DefaultBackendsDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultBackendsDir, "dev.tf"), []byte{}, 0644))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultBackendsDir, "readme.md"), []byte{}, 0644))

		got := List()

		assert.Equal(t, []string{"dev"}, got)
	})
}
