package fs

import "os"

type OsFs struct{}

func (o OsFs) ChangeWorkingDirectory(path string) error {
	return os.Chdir(path)
}
