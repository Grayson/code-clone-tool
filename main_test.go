package main

import (
	"grayson/cct/lib"
	git "grayson/cct/lib/GitApi"
	githubapi "grayson/cct/lib/GithubApi"
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
		// TODO: Add test cases.
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
		gc     git.GitApi
	}
	tests := []struct {
		name    string
		args    args
		want    lib.Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := performGitActions(tt.args.action, tt.args.gc)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
			if got := loadEnv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeEnvs(t *testing.T) {
	type args struct {
		change *lib.Env
		into   *lib.Env
	}
	tests := []struct {
		name string
		args args
		want *lib.Env
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeEnvs(tt.args.change, tt.args.into); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeEnvs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapActions(t *testing.T) {
	type args struct {
		repos *githubapi.GithubOrgReposResponse
	}
	tests := []struct {
		name        string
		args        args
		wantActions []lib.Action
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActions, err := mapActions(tt.args.repos)
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

func Test_discernPathInfo(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name  string
		args  args
		want  lib.PathExistential
		want1 lib.PathType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := discernPathInfo(tt.args.path)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("discernPathInfo() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("discernPathInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
