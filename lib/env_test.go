package lib_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/grayson/code-clone-tool/lib"
)

func Test_loadEnvironmentVariables(t *testing.T) {
	genget := func(m map[string]string) func(string) (string, bool) {
		return func(k string) (string, bool) {
			v, ok := m[k]
			return v, ok
		}
	}

	type args struct {
		get lib.GetEnvVar
	}
	tests := []struct {
		name string
		args args
		want *lib.Env
	}{
		{
			"Load all environment variables",
			args{genget(map[string]string{
				"PERSONAL_ACCESS_TOKEN": "pat",
				"API_URL":               "url",
				"WORKING_DIRECTORY":     "wd",
				"IS_MIRROR":             "true",
			}),
			},
			&lib.Env{
				PersonalAccessToken: "pat",
				ApiUrl:              "url",
				WorkingDirectory:    "wd",
				IsMirror:            "true",
			},
		},
		{
			"Load some environment variables",
			args{genget(map[string]string{
				"PERSONAL_ACCESS_TOKEN": "pat",
				"API_URL":               "url",
			}),
			},
			&lib.Env{
				PersonalAccessToken: "pat",
				ApiUrl:              "url",
			},
		},
		{
			"Load no environment variables",
			args{genget(map[string]string{})},
			&lib.Env{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.LoadEnvironmentVariables(tt.args.get); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadEnvironmentVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_EnsureTags(t *testing.T) {
	env := lib.Env{}
	rtype := reflect.TypeOf(env)
	fieldCount := rtype.NumField()
	for fieldIndex := 0; fieldIndex < fieldCount; fieldIndex++ {
		field := rtype.Field(fieldIndex)
		env := field.Tag.Get("env")
		if len(env) == 0 {
			t.Errorf("field %v of Env is missing the 'env' tag", field.Name)
		}
	}
}

func TestEnv_Merge(t *testing.T) {
	left := lib.Env{
		ApiUrl:              "apiurl",
		PersonalAccessToken: "pat",
		WorkingDirectory:    "wd",
	}
	right := lib.Env{
		ApiUrl:              "rurl",
		PersonalAccessToken: "rat",
		WorkingDirectory:    "rd",
	}
	empty := lib.Env{
		ApiUrl:              "",
		PersonalAccessToken: "",
		WorkingDirectory:    "",
	}

	type args struct {
		left  *lib.Env
		right *lib.Env
	}
	tests := []struct {
		name string
		args args
		want *lib.Env
	}{
		{
			"Choose left (nil right)",
			args{&left, nil},
			&left,
		},
		{
			"Choose left (empty right)",
			args{&left, &empty},
			&left,
		},
		{
			"Choose right (nil left)",
			args{nil, &right},
			&right,
		},
		{
			"Choose right (empty left)",
			args{&empty, &right},
			&right,
		},
		{
			"Choose right (full left)",
			args{&left, &right},
			&right,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.left.Merge(tt.args.right); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadEnvironmentYamlFile(t *testing.T) {
	y := func(yaml string) func() ([]byte, error) {
		return func() ([]byte, error) {
			if yaml == "err" {
				return nil, errors.New("expected")
			}
			return []byte(yaml), nil
		}
	}
	type args struct {
		read lib.ReadYamlFile
	}
	tests := []struct {
		name       string
		args       args
		want       *lib.Env
		expectsErr bool
	}{
		{
			"Load all items",
			args{y(`---
personal_access_token: pat
api_url: url
working_directory: wd
is_mirror: true`)},
			&lib.Env{
				PersonalAccessToken: "pat",
				ApiUrl:              "url",
				WorkingDirectory:    "wd",
				IsMirror:            "true",
			},
			false,
		},
		{
			"Load some items",
			args{y(`---
personal_access_token: pat
api_url: url`)},
			&lib.Env{
				PersonalAccessToken: "pat",
				ApiUrl:              "url",
			},
			false,
		},
		{
			"Load no items",
			args{y("")},
			&lib.Env{},
			false,
		},
		{
			"Load invalid yaml",
			args{y(`---
personal_access_token: pat
api_url: ,err`)},
			nil,
			true,
		},
		{
			"Handle read errors",
			args{y("err")},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lib.LoadEnvironmentYamlFile(tt.args.read)
			receivedExpectedErr := tt.expectsErr == (got == nil)
			if !(receivedExpectedErr || reflect.DeepEqual(got, tt.want)) {
				t.Errorf("LoadEnvironmentYamlFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
