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
	for _, repo := range *repos {
		log.Printf("Repo: %v @ %v", repo.FullName, repo.GitUrl)
	}
}
