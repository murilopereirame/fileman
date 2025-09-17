package fs

import "os"

type FileSystem interface {
	ReadDir(path string) ([]os.DirEntry, error)
	DeleteFile(path string) error
	ReadFile(path string) ([]byte, error)
}

type FS struct{}

// ReadDir reads a given path and returns its entries
func (f FS) ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

// ReadFile reads a given file and returns its entries
func (f FS) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// DeleteFile deletes a given file from the filesystem.
// If something went wrong, an error of type *PathError is returned.
func (f FS) DeleteFile(path string) error {
	return os.Remove(path)
}
