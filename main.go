package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"grayson/cct/lib"
	githubclient "grayson/cct/lib/GitApi"
	githubapi "grayson/cct/lib/GithubApi"
)

func loadEnv() *lib.Env {
	readers := []lib.ReadYamlFile{
		func() ([]byte, error) {
			return os.ReadFile(".env")
		},
	}
	return lib.NewEnv(os.LookupEnv, readers)
}

func main() {
	env := loadEnv()

	client := githubapi.NewClient(http.DefaultClient, env.PersonalAccessToken)
	resp, err := client.FetchOrgInformation(env.OrganizationUrl)

	if err != nil {
		log.Fatal(err)
	}

	if errResp, ok := resp.GetRight(); ok {
		log.Printf("Service error with the following message:\n%v\n\n%v", errResp.Message, errResp.DocumentationURL)
		return
	}

	repos, _ := resp.GetLeft()
	gc := githubclient.CreateGitClient(log.Default())
	for _, action := range mapActions(repos) {
		var output string
		var err error
		switch action.Task {
		case lib.Clone:
			output, err = gc.Clone(action.GitUrl, action.Path)
		case lib.Pull:
			output, err = gc.Pull(action.Path)
		default:
			panic(fmt.Sprintf("Unexpected task: %v", action.Task.String()))
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Print(output)
	}
}

func mapActions(repos *githubapi.GithubOrgReposResponse) (actions []lib.Action) {
	for _, repo := range *repos {
		action := lib.Action{
			Task:   lib.DiscernTask(repo.FullName, discernPathInfo),
			Path:   repo.FullName,
			GitUrl: repo.SshUrl,
		}
		actions = append(actions, action)
	}
	return
}

func discernPathInfo(path string) (lib.PathExistential, lib.PathType) {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return lib.DoesNotExist, lib.None
	}

	if info.IsDir() {
		return lib.Exists, lib.IsDirectory
	}
	return lib.Exists, lib.IsFile
}
