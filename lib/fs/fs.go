package fs

type PathExistential int

const (
	Exists PathExistential = iota
	DoesNotExist
)

type PathType int

const (
	None PathType = iota
	IsFile
	IsDirectory
)

type Fs interface {
	ChangeWorkingDirectory(string) error
	Info(path string) (PathExistential, PathType)
}
