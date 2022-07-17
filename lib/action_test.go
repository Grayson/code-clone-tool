package lib

import (
	"testing"
)

func TestDiscernTask(t *testing.T) {
	type args struct {
		infoDiscerner DiscernPathInfo
	}
	tests := []struct {
		name string
		args args
		want Task
	}{
		{
			"Is Directory that Exists",
			args{func(string) (PathExistential, PathType) { return Exists, IsDirectory }},
			Pull,
		},
		{
			"Does Not Exist",
			args{func(string) (PathExistential, PathType) { return DoesNotExist, IsDirectory }},
			Clone,
		},
		{
			"Is File",
			args{func(string) (PathExistential, PathType) { return Exists, IsFile }},
			Invalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiscernTask("test/path", tt.args.infoDiscerner); got != tt.want {
				t.Errorf("DiscernTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
