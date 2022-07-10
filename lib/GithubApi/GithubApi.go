package githubapi

import "grayson/cct/lib/either"

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
}

type GithubOrgReposResponse []GithubOrgReposResponseItem

type GithubApi interface {
	FetchOrgInformation() (*either.Either[*GithubOrgReposResponse, *GithubOrgReposErrorResponse], error)
}
