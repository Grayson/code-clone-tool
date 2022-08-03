package githubapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/grayson/code-clone-tool/lib/either"
)

const (
	pageLimit       = 32
	defaultPageSize = 30
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

func (c *GithubClient) FetchOrgInformation(urlString string) (out *either.Either[*GithubOrgReposResponse, *GithubOrgReposErrorResponse], err error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return
	}

	responses := make(GithubOrgReposResponse, 0)

	for page := 0; page < pageLimit; page++ {
		either, innerErr := getRepos(*u, page, c.personalAccessToken, *c.client)
		if innerErr != nil {
			err = innerErr
			return
		}
		if _, ok := either.GetRight(); ok {
			out = either
			return
		}

		additionalResponses, ok := either.GetLeft()
		if !ok {
			panic("unexpected case where we have neither responses nor errors from Github!")
		}
		responses = append(responses, *additionalResponses)

		if len(*additionalResponses) < pageLimit {
			break
		}
	}

	out = either.Of[*GithubOrgReposResponse, *GithubOrgReposErrorResponse](&responses)
	return
}

func getRepos(url url.URL, page int, pat string, client http.Client) (out *either.Either[*GithubOrgReposResponse, *GithubOrgReposErrorResponse], err error) {
	query := url.Query()
	query.Set("page", strconv.Itoa(page))
	url.RawQuery = query.Encode()

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
