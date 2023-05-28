package fs

type PathExistential int

const (
	Exists PathExistential = iota
	DoesNotExist
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=PathType
type PathType int

const (
	None PathType = iota
	IsFile
	IsDirectory
)

type Fs interface {
	ChangeWorkingDirectory(string) error
	GetWorkingDirectory() (string, error)
	Info(path string) (PathExistential, PathType)
}
