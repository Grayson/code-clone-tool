package lib

import "grayson/cct/lib/optional"

type Task int

const (
	Unknown Task = iota
	Invalid
	Clone
	Pull
)

type Action struct {
	task   Task
	path   string
	gitUrl string
}

func (act *Action) Execute() *optional.Optional[*error] {
	panic("Unimplemented")
}

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

type DiscernPathInfo func(path string) (PathExistential, PathType)

func DiscernTask(path string, infoDiscerner DiscernPathInfo) Task {
	existence, pathType := infoDiscerner(path)
	switch pathType {
	case None:
		return Invalid
	case IsFile:
		return Invalid
	case IsDirectory:
		switch existence {
		case Exists:
			return Pull
		case DoesNotExist:
			return Clone
		}
	}
	return Invalid
}
