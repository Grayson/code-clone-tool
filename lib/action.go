package lib

import (
	"github.com/grayson/code-clone-tool/lib/fs"
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=Task
type Task int

const (
	Unknown Task = iota
	Invalid
	Clone
	Pull
)

type Action struct {
	Task   Task
	Path   string
	GitUrl string
}

type DiscernPathInfo func(path string) (fs.PathExistential, fs.PathType)

func DiscernTask(path string, fsimpl fs.Fs) Task {
	existence, pathType := fsimpl.Info(path)
	switch existence {
	case fs.DoesNotExist:
		return Clone
	case fs.Exists:
		if fs.IsDirectory == pathType {
			return Pull
		}
	}
	return Invalid
}
