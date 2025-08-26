package handler

type File struct {
	createdAt int64
	age       float64
	name      string
	path      string
	isDir     bool
	error     []string
}

func NewFile(createdAt int64, age float64, name string, path string, isDir bool, error []string) *File {
	return &File{
		createdAt: createdAt,
		age:       age,
		name:      name,
		path:      path,
		isDir:     isDir,
		error:     error,
	}
}
