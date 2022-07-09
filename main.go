package main

import (
	"log"
	"net/http"
	"os"

	"grayson/cct/lib"
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
	log.Printf("%#v", mapActions(repos))
}

func mapActions(repos *githubapi.GithubOrgReposResponse) (actions []lib.Action) {
	for _, repo := range *repos {
		action := lib.Action{
			Task:   lib.DiscernTask(repo.FullName, discernPathInfo),
			Path:   repo.FullName,
			GitUrl: repo.GitUrl,
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
