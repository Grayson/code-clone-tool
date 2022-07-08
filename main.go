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

	tmp, _ := resp.GetLeft()
	log.Printf("%#v\n", tmp)
}
