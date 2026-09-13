package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExists(t *testing.T) {
	t.Run("ReturnTrueWhenFileExists", func(t *testing.T) {
		var dir = t.TempDir()
		var path = filepath.Join(dir, "test.tf")
		require.NoError(t, os.WriteFile(path, []byte(""), 0644))

		assert.True(t, Exists(path))
	})

	t.Run("ReturnTrueWhenDirectoryExists", func(t *testing.T) {
		var dir = t.TempDir()

		assert.True(t, Exists(dir))
	})

	t.Run("ReturnFalseWhenPathDoesNotExist", func(t *testing.T) {
		assert.False(t, Exists("/nonexistent/path/test.tf"))
	})
}
