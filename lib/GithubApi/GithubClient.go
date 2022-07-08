package githubapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"grayson/cct/lib/either"
)

type GithubClient struct {
	client              *http.Client
	personalAccessToken string
}

func NewClient(client *http.Client, personalAccessToken string) *GithubClient {
	return &GithubClient{
		client,
		personalAccessToken,
	}
}

func (c *GithubClient) FetchOrgInformation(url string) (out *either.Either[*GithubOrgReposResponse, *GithubOrgReposErrorResponse], err error) {
	out = nil
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", c.personalAccessToken))

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var errorResponse GithubOrgReposErrorResponse
	json.Unmarshal(bytes, &errorResponse)
	if errorResponse.Message != "" {
		out = either.Of[*GithubOrgReposResponse, *GithubOrgReposErrorResponse](&errorResponse)
		return
	}

	var repoResponse GithubOrgReposResponse
	json.Unmarshal(bytes, &repoResponse)
	out = either.Of[*GithubOrgReposResponse, *GithubOrgReposErrorResponse](repoResponse)

	return
}
