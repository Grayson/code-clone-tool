package app

import (
	"reflect"
	"testing"

	"github.com/grayson/code-clone-tool/lib"
	githubapi "github.com/grayson/code-clone-tool/lib/GithubApi"
	"github.com/grayson/code-clone-tool/lib/fs"
)

func Test_mapActions(t *testing.T) {
	type args struct {
		fs    fs.Fs
		repos *githubapi.GithubOrgReposResponse
	}
	tests := []struct {
		name        string
		args        args
		wantActions []lib.Action
		wantErr     bool
	}{
		{
			"Test",
			args{
				&TestFs{
					FileInfo: map[string]TestFsInfo{
						"clone": {fs.DoesNotExist, fs.None},
						"pull":  {fs.Exists, fs.IsDirectory},
						"file":  {fs.Exists, fs.IsFile},
					},
				},
				&githubapi.GithubOrgReposResponse{
					{
						FullName: "clone",
						SshUrl:   "ssh",
					},
					{
						FullName: "pull",
						SshUrl:   "ssh",
					},
					{
						FullName: "file",
						SshUrl:   "ssh",
					},
				},
			},
			[]lib.Action{
				{Task: lib.Clone, Path: "clone", GitUrl: "ssh"},
				{Task: lib.Pull, Path: "pull", GitUrl: "ssh"},
				{Task: lib.Invalid, Path: "file", GitUrl: "ssh"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActions, err := mapActions(tt.args.fs, tt.args.repos)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotActions, tt.wantActions) {
				t.Errorf("mapActions() = %v, want %v", gotActions, tt.wantActions)
			}
		})
	}
}

// func Test_performGitActions(t *testing.T) {
// 	type args struct {
// 		action lib.Action
// 		gc     TestGit
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    lib.Task
// 		wantErr bool
// 	}{
// 		{
// 			"Test Clone",
// 			args{
// 				lib.Action{
// 					Task:   lib.Clone,
// 					Path:   "path",
// 					GitUrl: "ssh://git-repo",
// 				},
// 				TestGit{},
// 			},
// 			lib.Clone,
// 			false,
// 		},
// 		{
// 			"Test Pull",
// 			args{
// 				lib.Action{
// 					Task:   lib.Pull,
// 					Path:   "path",
// 					GitUrl: "ssh://git-repo",
// 				},
// 				TestGit{},
// 			},
// 			lib.Pull,
// 			false,
// 		},
// 		{
// 			"Test Invalid",
// 			args{
// 				lib.Action{
// 					Task:   lib.Invalid,
// 					Path:   "path",
// 					GitUrl: "ssh://git-repo",
// 				},
// 				TestGit{},
// 			},
// 			lib.Invalid,
// 			true,
// 		},
// 		{
// 			"Test Unknown",
// 			args{
// 				lib.Action{
// 					Task:   lib.Unknown,
// 					Path:   "path",
// 					GitUrl: "ssh://git-repo",
// 				},
// 				TestGit{},
// 			},
// 			lib.Invalid,
// 			true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := performGitActions(tt.args.action, 0, &tt.args.gc)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("performGitActions() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("performGitActions() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
