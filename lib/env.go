package lib

import (
	"gopkg.in/yaml.v3"
)

type Env struct {
	PersonalAccessToken string
}

type GetEnvVar func(string) (string, bool)
type ReadYamlFile func() ([]byte, error)

type LoadStatus int8

type envFile struct {
	PersonalAccessToken string `yaml:"personal_access_token"`
}

type envFileWrapper struct {
	EnvFile envFile
	Status  LoadStatus
}

type loadYamlEnvFile func() *envFileWrapper

const (
	Unloaded LoadStatus = iota
	Loaded
	FailedToLoad
)

func loadEnvFile(reader ReadYamlFile) *envFileWrapper {
	var file envFile
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
	value, ok := get("GH_PERSONAL_ACCESS_TOKEN")
	if ok {
		return value
	}
	if env := load(); env.Status == Loaded {
		return env.EnvFile.PersonalAccessToken
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
	}
}
