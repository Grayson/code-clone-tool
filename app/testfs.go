package app

import "github.com/grayson/code-clone-tool/lib/fs"

type TestFsInfo struct {
	E fs.PathExistential
	T fs.PathType
}

type TestFs struct {
	Error    error
	FileInfo map[string]TestFsInfo
}

func (f *TestFs) ChangeWorkingDirectory(_ string) error {
	return f.Error
}

func (*TestFs) GetWorkingDirectory() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (f *TestFs) Info(path string) (fs.PathExistential, fs.PathType) {
	x := f.FileInfo[path]
	return x.E, x.T
}
