package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashFile(t *testing.T) {
	t.Run("ReturnConsistentHash", func(t *testing.T) {
		var path = filepath.Join(t.TempDir(), "test.tf")
		require.NoError(t, os.WriteFile(path, []byte("hello"), 0644))

		hash1, err := HashFile(path)
		require.NoError(t, err)

		hash2, err := HashFile(path)
		require.NoError(t, err)

		assert.Equal(t, hash1, hash2)
		assert.Len(t, hash1, 64)
	})

	t.Run("ReturnDifferentHashForDifferentContent", func(t *testing.T) {
		dir := t.TempDir()
		path1 := filepath.Join(dir, "a.tf")
		path2 := filepath.Join(dir, "b.tf")
		require.NoError(t, os.WriteFile(path1, []byte("hello"), 0644))
		require.NoError(t, os.WriteFile(path2, []byte("world"), 0644))

		hash1, err := HashFile(path1)
		require.NoError(t, err)

		hash2, err := HashFile(path2)
		require.NoError(t, err)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("ErrorWhenFileDoesNotExist", func(t *testing.T) {
		_, err := HashFile("/nonexistent/file.tf")
		assert.Error(t, err)
	})
}

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
