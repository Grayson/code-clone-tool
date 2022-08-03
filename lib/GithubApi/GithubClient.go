package githubapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	neturl "net/url"

	"github.com/grayson/code-clone-tool/lib/either"
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
	u, err := neturl.Parse(url)
	if err != nil {
		return
	}
	return getRepos(*u, c.personalAccessToken, *c.client)
}

func getRepos(url url.URL, pat string, client http.Client) (out *either.Either[*GithubOrgReposResponse, *GithubOrgReposErrorResponse], err error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", pat))

	resp, err := client.Do(req)
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
	out = either.Of[*GithubOrgReposResponse, *GithubOrgReposErrorResponse](&repoResponse)
	return
}
