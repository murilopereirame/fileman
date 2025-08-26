package fs

import "os"

type FileSystem interface {
	ReadDir(path string) ([]os.DirEntry, error)
	DeleteFile(path string) error
}

// ReadDir reads a given path and returns its entries
func ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

// DeleteFile deletes a given file from the filesystem.
// If something went wrong, an error of type *PathError is returned.
func DeleteFile(path string) error {
	return os.Remove(path)
}
