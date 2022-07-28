package lib_test

import (
	"grayson/cct/lib"
	"reflect"
	"testing"
)

func TestEmptyPATWhenEnvironmentGetterReturnsEmptyValue(t *testing.T) {
	get := func(key string) (string, bool) {
		return "", false
	}
	env := lib.NewEnv(get, []lib.ReadYamlFile{})
	if env.PersonalAccessToken != "" {
		t.Errorf("Receive non-empty personal access token (%v)", env.PersonalAccessToken)
	}
}

func TestValidPATFromEnvVarGetterFunction(t *testing.T) {
	expect := "PAT"
	envMap := map[string]string{
		"PERSONAL_ACCESS_TOKEN": expect,
	}
	get := func(key string) (string, bool) {
		val, ok := envMap[key]
		return val, ok
	}

	env := lib.NewEnv(get, []lib.ReadYamlFile{})
	if env.PersonalAccessToken != expect {
		t.Errorf("Receive unexpected personal access token (received: %v, expected: %v)", env.PersonalAccessToken, expect)
	}
}

func TestValidPATHFromFallbackToEnvVarFile(t *testing.T) {
	const (
		expect = "PAT"
		yaml   = `---
personal_access_token: PAT`
	)
	get := func(key string) (string, bool) { return "", false }
	reader := func() ([]byte, error) {
		return []byte(yaml), nil
	}

	env := lib.NewEnv(get, []lib.ReadYamlFile{reader})
	if env.PersonalAccessToken != expect {
		t.Errorf("Receive unexpected personal access token (received: %v, expected: %v)", env.PersonalAccessToken, expect)
	}
}

func TestEnvVarOverridingEnvFile(t *testing.T) {
	const (
		expect = "PAT"
		yaml   = `---
personal_access_token: NOT_PAT`
	)
	envMap := map[string]string{"PERSONAL_ACCESS_TOKEN": expect}
	get := func(key string) (string, bool) {
		val, ok := envMap[key]
		return val, ok
	}
	reader := func() ([]byte, error) {
		return []byte(yaml), nil
	}

	env := lib.NewEnv(get, []lib.ReadYamlFile{reader})
	if env.PersonalAccessToken != expect {
		t.Errorf("Receive unexpected personal access token (received: %v, expected: %v)", env.PersonalAccessToken, expect)
	}
}

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
		return func() ([]byte, error) { return []byte(yaml), nil }
	}
	type args struct {
		read lib.ReadYamlFile
	}
	tests := []struct {
		name string
		args args
		want *lib.Env
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
		},
		{
			"Load no items",
			args{y("")},
			&lib.Env{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.LoadEnvironmentYamlFile(tt.args.read); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadEnvironmentYamlFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
