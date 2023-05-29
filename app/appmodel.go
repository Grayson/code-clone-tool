package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grayson/code-clone-tool/lib"
	githubapi "github.com/grayson/code-clone-tool/lib/GithubApi"
	"github.com/grayson/code-clone-tool/lib/fs"
)

type AppModel struct {
	env *lib.Env

	Error       error
	children    []tea.Model
	headerStyle lipgloss.Style
}

type errMsg error

type configurationCompleteMsg struct {
	personalAccessToken string
	url                 string
	isMirror            bool
}

type repoResponseMsg struct {
	repos *githubapi.GithubOrgReposResponse
}

func InitAppModel(env *lib.Env, version string, fileSystem fs.Fs) *AppModel {
	versionStyle := lipgloss.NewStyle().Faint(true).SetString(version)
	header := fmt.Sprintf("code-clone-tool \n%v", versionStyle)

	return &AppModel{
		env: env,
		children: []tea.Model{
			NewConfigModel(env, fileSystem),
			&fetchreposmodel{},
			NewPerformGitActionsModel(fileSystem),
			&counttasksmodel{},
		},
		headerStyle: lipgloss.NewStyle().
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			SetString(header),
	}
}

func (app *AppModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for _, m := range app.children {
		cmd := m.Init()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func (app *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch actual := msg.(type) {
	case tea.KeyMsg:
		cmd := handleKeyboardEvent(actual, app)
		if cmd != nil {
			return app, cmd
		}
	case errMsg:
		app.Error = actual
		return app, tea.Quit
	}

	cmds := make([]tea.Cmd, 0)
	for _, model := range app.children {
		_, cmd := model.Update(msg)
		cmds = append(cmds, cmd)
	}

	return app, tea.Batch(cmds...)
}

func (app *AppModel) View() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%v\n", app.headerStyle.String())
	if app.Error != nil {
		fmt.Fprintf(&sb, "%v", app.Error)
		return sb.String()
	}

	for _, m := range app.children {
		fmt.Fprintln(&sb, m.View())
	}

	return sb.String()
}

func handleKeyboardEvent(msg tea.KeyMsg, app *AppModel) tea.Cmd {
	switch msg.String() {
	case "q", tea.KeyCtrlC.String(), tea.KeyEsc.String():
		return tea.Quit
	}
	return nil
}

func reportError(err error) tea.Cmd {
	return func() tea.Msg {
		return errMsg(err)
	}
}
