package lib

import (
	"testing"

	"github.com/grayson/code-clone-tool/lib/fs"
)

func TestDiscernTask(t *testing.T) {
	testfs := TestFs{
		map[string]fs_info{
			"dne": {
				fs.DoesNotExist,
				fs.None,
			},
			"eid": {
				fs.Exists,
				fs.IsDirectory,
			},
			"eif": {
				fs.Exists,
				fs.IsFile,
			},
		},
	}

	type args struct {
		path   string
		fsimpl fs.Fs
	}
	tests := []struct {
		name string
		args args
		want Task
	}{
		{
			"Does Not Exist -> Clone",
			args{"dne", &testfs},
			Clone,
		},
		{
			"Exists, Is Directory -> Pull",
			args{"eid", &testfs},
			Pull,
		},
		{
			"Exists, Is File -> Invalid",
			args{"eif", &testfs},
			Invalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiscernTask(tt.args.path, tt.args.fsimpl); got != tt.want {
				t.Errorf("DiscernTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

type fs_info struct {
	Existential fs.PathExistential
	Type        fs.PathType
}

type TestFs struct {
	filesytem map[string]fs_info
}

func (fs *TestFs) ChangeWorkingDirectory(_ string) error {
	panic("not implemented") // TODO: Implement
}

func (fs *TestFs) GetWorkingDirectory() (string, error) {
	panic("not implemented")
}

func (fs *TestFs) Info(path string) (fs.PathExistential, fs.PathType) {
	out := fs.filesytem[path]
	return out.Existential, out.Type
}
