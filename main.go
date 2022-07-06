package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"grayson/cct/lib"
)

func loadEnv() *lib.Env {
	readers := []lib.ReadYamlFile{
		func() ([]byte, error) {
			return os.ReadFile(".env")
		},
	}
	return lib.NewEnv(os.LookupEnv, readers)
}

type GithubOrgReposErrorResponse struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

type GithubOrgReposResponseItem struct {
	Identifier int    `json:"id"`
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
}

type GithubOrgReposeResponse []GithubOrgReposResponseItem

func main() {
	env := loadEnv()

	req, err := http.NewRequest("GET", "https://api.github.com/orgs/<org-name>/repos", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", env.PersonalAccessToken))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tmp GithubOrgReposeResponse
	json.Unmarshal(bytes, &tmp)

	log.Println(tmp)
}
