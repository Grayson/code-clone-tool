package app

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grayson/code-clone-tool/lib"
	git "github.com/grayson/code-clone-tool/lib/GitApi"
	githubapi "github.com/grayson/code-clone-tool/lib/GithubApi"
	"github.com/grayson/code-clone-tool/lib/fs"
)

type performGitActionsModel struct {
	log          *log.Logger
	fileSystem   fs.Fs
	shouldMirror bool

	progress  progress.Model
	total     int
	completed int

	actions []lib.Action
	api     git.GitApi
	state   performingGitActionsState
}

type doingGitActionMsg struct {
	index int
}
type finishedGitActionMsg int
type finishedPerformingGitActions struct{}

type performingGitActionsState int

const (
	waitingToPerformGitActionsState performingGitActionsState = iota
	updatingPerformingGitActionsState
	finishedPerformingGitActionsState
)

func NewPerformGitActionsModel(fileSystem fs.Fs) *performGitActionsModel {
	file, _ := tea.LogToFile("git_actions.log", "debug")
	log := log.New(file, "", 0)
	// TODO: Cleanup file

	progress := progress.New()

	return &performGitActionsModel{
		fileSystem: fileSystem,
		log:        log,
		progress:   progress,
	}
}

func (m *performGitActionsModel) Init() tea.Cmd {
	return nil
}

func (m *performGitActionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch actual := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = actual.Width - 8
		return m, nil
	case repoResponseMsg:
		actions, err := mapActions(m.fileSystem, actual.repos)
		if err != nil {
			return m, reportError(err)
		}

		m.api = determineGitClient(m.shouldMirror, m.log)
		m.actions = actions
		m.total = len(actions)
		m.state = updatingPerformingGitActionsState

		return m, startGitActions(m)
	case configurationCompleteMsg:
		m.shouldMirror = actual.isMirror
	case doingGitActionMsg:
		return m, performGitAction(m.actions, actual.index, m.api)
	case finishedGitActionMsg:
		m.completed++
		if m.completed == m.total {
			m.state = finishedPerformingGitActionsState
			return m, func() tea.Msg { return finishedPerformingGitActions{} }
		}
		return m, performGitAction(m.actions, int(actual), m.api)
	}

	return m, nil
}

func (m *performGitActionsModel) View() string {
	if m.state == waitingToPerformGitActionsState {
		return ""
	}
	return fmt.Sprintf("%v of %v finished\n%v", m.completed, m.total, m.progress.ViewAs(float64(m.completed)/float64(m.total)))
}

func startGitActions(m *performGitActionsModel) tea.Cmd {
	return func() tea.Msg {
		return doingGitActionMsg{0}
	}
}

func mapActions(fs fs.Fs, repos *githubapi.GithubOrgReposResponse) (actions []lib.Action, err error) {
	for _, repo := range *repos {
		action := lib.Action{
			Task:   lib.DiscernTask(repo.FullName, fs),
			Path:   repo.FullName,
			GitUrl: repo.SshUrl,
		}
		actions = append(actions, action)
	}
	return
}

func performGitAction(actions []lib.Action, index int, api git.GitApi) tea.Cmd {
	return func() tea.Msg {
		var err error
		action := actions[index]
		// task := action.Task
		switch action.Task {
		case lib.Clone:
			_, err = api.Clone(action.GitUrl, action.Path)
		case lib.Pull:
			_, err = api.Pull(action.Path)
		default:
			err = fmt.Errorf("unexpected task: %v", action.Task.String())
		}
		if err != nil {
			return reportError(err)
		}
		// TODO: print result of git command to file, see `_` usages
		return finishedGitActionMsg(index + 1)
	}
}

func determineGitClient(isMirror bool, log *log.Logger) git.GitApi {
	if isMirror {
		return git.CreateMirrorClient(log)
	}
	return git.CreateGitClient(log)
}
