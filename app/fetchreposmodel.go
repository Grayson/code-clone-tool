package app

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	githubapi "github.com/grayson/code-clone-tool/lib/GithubApi"
)

type fetchreposstate int

const (
	waitingForConfiguration fetchreposstate = iota
	fetching
	done
)

type fetchreposmodel struct {
	spinner    spinner.Model
	found      int
	state      fetchreposstate
	willMirror bool
}

func (m *fetchreposmodel) Init() tea.Cmd {
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot
	return nil
}

func (m *fetchreposmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch actual := msg.(type) {
	case configurationCompleteMsg:
		m.state = fetching
		m.willMirror = actual.isMirror
		return m, fetchRepoInformationCmd(actual.personalAccessToken, actual.url)
	case repoResponseMsg:
		m.state = done
		m.found = len(*actual.repos)
	}

	if m.state == fetching {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *fetchreposmodel) View() string {
	if m.state == waitingForConfiguration {
		return ""
	}

	if m.state == fetching {
		return fmt.Sprintf("%v Fetching repo information...", m.spinner.View())
	}

	verb := "Cloning"
	if m.willMirror {
		verb = "Mirroring"
	}

	return fmt.Sprintf("%v %v repos", verb, m.found)
}

func fetchRepoInformationCmd(pat string, url string) tea.Cmd {
	return func() tea.Msg {
		client := githubapi.NewClient(http.DefaultClient, pat)

		resp, err := client.FetchOrgInformation(url)
		if err != nil {
			return errMsg(err)
		}

		if errResp, ok := resp.GetRight(); ok {
			return errMsg(fmt.Errorf("service error with the following message:\n%v\n\n%v", errResp.Message, errResp.DocumentationURL))
		}

		repos, _ := resp.GetLeft()
		return repoResponseMsg{repos}
	}
}
