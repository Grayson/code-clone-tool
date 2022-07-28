package lib

import (
	"fmt"
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

type envFileWrapper struct {
	EnvFile EnvFile
	Status  LoadStatus
}

type loadYamlEnvFile func() *envFileWrapper

const (
	Unloaded LoadStatus = iota
	Loaded
	FailedToLoad
)

func LoadEnvironmentVariables(get GetEnvVar) *Env {
	env := Env{}
	rtype := reflect.TypeOf(env)
	relem := reflect.ValueOf(&env).Elem()
	fieldCount := rtype.NumField()
	for fieldIndex := 0; fieldIndex < fieldCount; fieldIndex++ {
		field := rtype.Field(fieldIndex)
		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}
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

func (e *envFileWrapper) lookup(key EnvironmentVariableKey) string {
	switch key {
	case PersonalAccessToken:
		return e.EnvFile.PersonalAccessToken
	case ApiUrl:
		return e.EnvFile.ApiUrl
	case WorkingDirectory:
		return e.EnvFile.WorkingDirectory
	}
	panic(fmt.Sprintf("Unknown key: %v", key))
}

func loadEnvFile(reader ReadYamlFile) *envFileWrapper {
	var file EnvFile
	bytes, err := reader()
	if err != nil {
		return &envFileWrapper{
			Status: FailedToLoad,
		}
	}
	if err := yaml.Unmarshal(bytes, &file); err != nil {
		return &envFileWrapper{
			Status: FailedToLoad,
		}
	}
	return &envFileWrapper{
		EnvFile: file,
		Status:  Loaded,
	}
}

func getPersonalAccessToken(get GetEnvVar, load loadYamlEnvFile) string {
	value, ok := get(string(PersonalAccessToken))
	if ok {
		return value
	}
	if env := load(); env.Status == Loaded {
		return env.EnvFile.PersonalAccessToken
	}

	return ""
}

func getOrganizationUrl(get GetEnvVar, load loadYamlEnvFile) string {
	value, ok := get(string(ApiUrl))
	if ok {
		return value
	}
	if env := load(); env.Status == Loaded {
		return env.EnvFile.ApiUrl
	}

	return ""
}

func getStringConfig(get GetEnvVar, load loadYamlEnvFile, key EnvironmentVariableKey) string {
	value, ok := get(string(key))
	if ok {
		return value
	}
	if env := load(); env.Status == Loaded {
		return env.lookup(key)
	}

	return ""
}

func NewEnv(get GetEnvVar, read []ReadYamlFile) *Env {
	envFile := &envFileWrapper{}
	loader := func() *envFileWrapper {
		if envFile.Status == Loaded {
			return envFile
		}

		for _, x := range read {
			tmp := loadEnvFile(x)
			if tmp.Status != Loaded {
				continue
			}
			envFile = tmp
			break
		}
		return envFile
	}

	return &Env{
		PersonalAccessToken: getPersonalAccessToken(get, loader),
		ApiUrl:              getOrganizationUrl(get, loader),
		WorkingDirectory:    getStringConfig(get, loader, WorkingDirectory),
	}
}
