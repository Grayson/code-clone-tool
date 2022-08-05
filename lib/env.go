package lib

import (
	"reflect"

	"gopkg.in/yaml.v3"
)

type Env struct {
	PersonalAccessToken string     `env:"PERSONAL_ACCESS_TOKEN" yaml:"personal_access_token"`
	ApiUrl              string     `env:"API_URL" yaml:"api_url"`
	WorkingDirectory    string     `env:"WORKING_DIRECTORY" yaml:"working_directory"`
	IsMirror            BoolString `env:"IS_MIRROR" yaml:"is_mirror"`
}

type GetEnvVar func(string) (string, bool)
type ReadYamlFile func() ([]byte, error)

type LoadStatus int8

type EnvironmentVariableKey string

const (
	PersonalAccessToken EnvironmentVariableKey = "PERSONAL_ACCESS_TOKEN"
	ApiUrl              EnvironmentVariableKey = "API_URL"
	WorkingDirectory    EnvironmentVariableKey = "WORKING_DIRECTORY"
)

type EnvFile struct {
	PersonalAccessToken string `yaml:"personal_access_token"`
	ApiUrl              string `yaml:"api_url"`
	WorkingDirectory    string `yaml:"working_directory"`
}

func LoadEnvironmentVariables(get GetEnvVar) *Env {
	env := Env{}
	rtype := reflect.TypeOf(env)
	relem := reflect.ValueOf(&env).Elem()
	fieldCount := rtype.NumField()
	for fieldIndex := 0; fieldIndex < fieldCount; fieldIndex++ {
		field := rtype.Field(fieldIndex)
		if x, ok := get(field.Tag.Get("env")); ok {
			relem.FieldByIndex(field.Index).SetString(x)
		}
	}
	return &env
}

func LoadEnvironmentYamlFile(read ReadYamlFile) *Env {
	var env Env
	bytes, err := read()
	if err != nil {
		return nil
	}
	if err := yaml.Unmarshal(bytes, &env); err != nil {
		return nil
	}
	return &env
}

func (left *Env) Merge(right *Env) *Env {
	if left == nil {
		return right
	}

	combined := *left

	if right == nil {
		return &combined
	}

	rvalue := reflect.ValueOf(right).Elem()
	rcombined := reflect.ValueOf(&combined).Elem()
	fieldCount := rvalue.NumField()
	for fieldIndex := 0; fieldIndex < fieldCount; fieldIndex++ {
		field := rvalue.Field(fieldIndex)
		if value := field.String(); value != "" {
			rcombined.Field(fieldIndex).SetString(value)
		}
	}
	return &combined
}
