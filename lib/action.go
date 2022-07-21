package lib

import (
	"fmt"
	"grayson/cct/lib/fs"
	"grayson/cct/lib/optional"
)

type Task int

const (
	Unknown Task = iota
	Invalid
	Clone
	Pull
)

func (t Task) String() string {
	switch t {
	case Unknown:
		return "Unknown"
	case Invalid:
		return "Invalid"
	case Clone:
		return "Clone"
	case Pull:
		return "Pull"
	}
	panic(fmt.Sprintf("Unexpected task case %s", t.String()))
}

type Action struct {
	Task   Task
	Path   string
	GitUrl string
}

func (act *Action) Execute() *optional.Optional[*error] {
	panic("Unimplemented")
}

type DiscernPathInfo func(path string) (fs.PathExistential, fs.PathType)

func DiscernTask(path string, fs fs.Fs) Task {
	existence, pathType := fs.Info(path)
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
