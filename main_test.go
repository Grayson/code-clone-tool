package main

import (
	"fmt"
	"grayson/cct/lib"
	githubapi "grayson/cct/lib/GithubApi"
	"grayson/cct/lib/either"
	"grayson/cct/lib/fs"
	"reflect"
	"testing"
)

func Test_countTasks(t *testing.T) {
	type args struct {
		tasks []lib.Task
	}
	tests := []struct {
		name    string
		args    args
		want    map[lib.Task]int
		wantErr bool
	}{
		{
			"Test 1 clone",
			args{[]lib.Task{lib.Clone}},
			map[lib.Task]int{lib.Clone: 1},
			false,
		},
		{
			"Test 1 pull",
			args{[]lib.Task{lib.Pull}},
			map[lib.Task]int{lib.Pull: 1},
			false,
		},
		{
			"Test 1 pull and 1 clone",
			args{[]lib.Task{lib.Pull, lib.Clone}},
			map[lib.Task]int{lib.Pull: 1, lib.Clone: 1},
			false,
		},
		{
			"Test 1 pull and 2 clone",
			args{[]lib.Task{lib.Pull, lib.Clone, lib.Clone}},
			map[lib.Task]int{lib.Pull: 1, lib.Clone: 2},
			false,
		},
		{
			"Test invalid case",
			args{[]lib.Task{lib.Invalid}},
			map[lib.Task]int{lib.Invalid: 1},
			false,
		},
		{
			"Test 1 unknown case",
			args{[]lib.Task{lib.Unknown}},
			map[lib.Task]int{lib.Unknown: 1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countTasks(tt.args.tasks)
			if (err != nil) != tt.wantErr {
				t.Errorf("countTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_performGitActions(t *testing.T) {
	type args struct {
		action lib.Action
		gc     TestGit
	}
	tests := []struct {
		name    string
		args    args
		want    lib.Task
		wantErr bool
	}{
		{
			"Test Clone",
			args{
				lib.Action{
					Task:   lib.Clone,
					Path:   "path",
					GitUrl: "ssh://git-repo",
				},
				TestGit{},
			},
			lib.Clone,
			false,
		},
		{
			"Test Pull",
			args{
				lib.Action{
					Task:   lib.Pull,
					Path:   "path",
					GitUrl: "ssh://git-repo",
				},
				TestGit{},
			},
			lib.Pull,
			false,
		},
		{
			"Test Invalid",
			args{
				lib.Action{
					Task:   lib.Invalid,
					Path:   "path",
					GitUrl: "ssh://git-repo",
				},
				TestGit{},
			},
			lib.Invalid,
			true,
		},
		{
			"Test Unknown",
			args{
				lib.Action{
					Task:   lib.Unknown,
					Path:   "path",
					GitUrl: "ssh://git-repo",
				},
				TestGit{},
			},
			lib.Invalid,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := performGitActions(tt.args.action, &tt.args.gc)
			if (err != nil) != tt.wantErr {
				t.Errorf("performGitActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("performGitActions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fetchRepoInformation(t *testing.T) {
	success := &TestApi{Response: &githubapi.GithubOrgReposResponse{}}
	type args struct {
		client githubapi.GithubApi
		url    string
	}
	tests := []struct {
		name    string
		args    args
		want    *githubapi.GithubOrgReposResponse
		wantErr bool
	}{
		{
			"Produces error result in case of local error",
			args{&TestApi{Error: fmt.Errorf("err")}, "url"},
			nil,
			true,
		},
		{
			"Produces error in case of Github Error Message",
			args{&TestApi{ErrorResponse: &githubapi.GithubOrgReposErrorResponse{}}, "url"},
			nil,
			true,
		},
		{
			"Produces value in case of success",
			args{success, "url"},
			success.Response,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchRepoInformation(tt.args.client, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchRepoInformation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchRepoInformation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cwd(t *testing.T) {
	type args struct {
		f fs.Fs
		p string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"No error",
			args{&TestFs{}, "path"},
			true,
			false,
		},
		{
			"Error",
			args{&TestFs{Error: fmt.Errorf("err")}, "path"},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cwd(tt.args.f, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("cwd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadEnv(t *testing.T) {
	tests := []struct {
		name string
		want *lib.Env
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadEnv(".env"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func Test_determineConfigPath(t *testing.T) {
	type args struct {
		initial  string
		fallback func() (string, bool)
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty initial and no env var lookup",
			args: args{"", func() (string, bool) { return "", false }},
			want: ".env",
		},
		{
			name: "Filled initial and no env var lookup",
			args: args{"initial", func() (string, bool) { return "", false }},
			want: "initial",
		},
		{
			name: "Filled initial and valid env var lookup",
			args: args{"initial", func() (string, bool) { return "env", true }},
			want: "initial",
		},
		{
			name: "Empty initial and valid env var lookup",
			args: args{"", func() (string, bool) { return "env", true }},
			want: "env",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineConfigPath(tt.args.initial, tt.args.fallback); got != tt.want {
				t.Errorf("determineConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test implementations

type TestGit struct {
	DidClone bool
	DidPull  bool
}

func (g *TestGit) Clone(gitUrl string, path string) (string, error) {
	g.DidClone = true
	return "stdout", nil
}

func (g *TestGit) Pull(destinationDir string) (string, error) {
	g.DidPull = true
	return "stdout", nil
}

type TestApi struct {
	Error         error
	ErrorResponse *githubapi.GithubOrgReposErrorResponse
	Response      *githubapi.GithubOrgReposResponse
}

func (api *TestApi) FetchOrgInformation(url string) (*either.Either[*githubapi.GithubOrgReposResponse, *githubapi.GithubOrgReposErrorResponse], error) {
	of := either.Of[*githubapi.GithubOrgReposResponse, *githubapi.GithubOrgReposErrorResponse]

	if api.Error != nil {
		return nil, api.Error
	}
	if api.ErrorResponse != nil {
		return of(api.ErrorResponse), nil
	}
	return of(api.Response), nil
}

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

func (f *TestFs) Info(path string) (fs.PathExistential, fs.PathType) {
	x := f.FileInfo[path]
	return x.E, x.T
}
