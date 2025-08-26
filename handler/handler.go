package handler

import (
	"container/list"
	"fileman/clock"
	"fileman/fs"
	"path/filepath"
)

type IFileHandler interface {
	ListFiles(fs fs.FileSystem, path string) (list.List, error)
}

type FileHandler struct {
	clock clock.Clock
}

func New(clock clock.Clock) *FileHandler {
	return &FileHandler{
		clock: clock,
	}
}

// ListFiles list the files in a given directory returning
// the files together with its details
func (f FileHandler) ListFiles(fs fs.FileSystem, path string) (list.List, error) {
	files := list.List{}
	dirEntries, err := fs.ReadDir(path)

	if err != nil {
		return list.List{}, err
	}

	for _, entry := range dirEntries {
		file := &File{
			name: entry.Name(),
		}

		info, err := entry.Info()

		if err != nil {
			file.error = err
		} else {
			file.createdAt = info.ModTime().Unix()
			file.age = f.clock.CalculateAge(info.ModTime().Unix())
			file.path = filepath.Join(path, entry.Name())
			file.isDir = info.IsDir()
		}

		files.PushBack(file)
	}

	return files, nil
}

// DeleteOldFiles deletes files older than the given threshold (in days)
// from the given path. It returns a list of errors encountered during the process.
func (f FileHandler) DeleteOldFiles(fs fs.FileSystem, path string, threshold float64) []error {
	files, err := f.ListFiles(fs, path)
	errors := make([]error, 0)

	if err != nil {
		errors = append(errors, err)
		return errors
	}

	for e := files.Front(); e != nil; e = e.Next() {
		file := e.Value.(*File)

		if file.error != nil {
			errors = append(errors, file.error)
			continue
		}

		if !file.isDir && file.age > threshold {
			err := fs.DeleteFile(file.path)

			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	return errors
}
