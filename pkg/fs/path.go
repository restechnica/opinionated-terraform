package fs

import "os"

// Exists checks whether a path exists on the filesystem.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
