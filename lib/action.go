package lib

import (
	"fmt"
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
		switch existence {
		case Exists:
			return Invalid
		case DoesNotExist:
			return Clone
		}
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
