package fs

import (
	"os"
)

type OsFs struct{}

func (o OsFs) ChangeWorkingDirectory(path string) error {
	return os.Chdir(path)
}

func (o OsFs) Info(path string) (PathExistential, PathType) {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return DoesNotExist, None
	}

	if info.IsDir() {
		return Exists, IsDirectory
	}
	return Exists, IsFile
}
