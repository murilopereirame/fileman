package fs

import "os"

type FileSystem interface {
	ReadDir(path string) ([]os.DirEntry, error)
}

func ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}
