package lib_test

import (
	"grayson/cct/lib"
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
		"GH_PERSONAL_ACCESS_TOKEN": expect,
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
	envMap := map[string]string{"GH_PERSONAL_ACCESS_TOKEN": expect}
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
