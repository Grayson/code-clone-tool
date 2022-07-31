package githubapi

import "github.com/grayson/code-clone-tool/lib/either"

type GithubOrgReposErrorResponse struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

type GithubOrgReposResponseItem struct {
	Identifier int    `json:"id"`
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
	GitUrl     string `json:"git_url"`
	HtmlUrl    string `json:"html_url"`
	SshUrl     string `json:"ssh_url"`
}

type GithubOrgReposResponse []GithubOrgReposResponseItem

type GithubApi interface {
	FetchOrgInformation(url string) (*either.Either[*GithubOrgReposResponse, *GithubOrgReposErrorResponse], error)
}
