package handler

import (
	"container/list"
	"fileman/clock"
	"fileman/fs"
	"path/filepath"
)

type IFileHandler interface {
	ListFiles(fs fs.FileSystem, path string) list.List
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
func (f FileHandler) ListFiles(fs fs.FileSystem, path string) list.List {
	files := list.List{}
	dirEntries, err := fs.ReadDir(path)

	if err != nil {
		return list.List{}
	}

	for _, entry := range dirEntries {
		file := &File{
			error: make([]string, 0),
		}

		info, err := entry.Info()

		if err != nil {
			file.error = append(file.error, err.Error())
		} else {
			file.name = entry.Name()
			file.createdAt = info.ModTime().Unix()
			file.age = f.clock.CalculateAge(info.ModTime().Unix())
			file.path = filepath.Join(path, entry.Name())
			file.isDir = info.IsDir()
		}

		files.PushBack(file)
	}

	return files
}
