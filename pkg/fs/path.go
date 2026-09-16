package fs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

// Exists checks whether a path exists on the filesystem.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// HashFile returns the SHA-256 hex digest of a file's contents.
func HashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("hashing file %q: %w", path, err)
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
