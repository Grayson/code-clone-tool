package app

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync/atomic"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grayson/code-clone-tool/lib"
	git "github.com/grayson/code-clone-tool/lib/GitApi"
	githubapi "github.com/grayson/code-clone-tool/lib/GithubApi"
	"github.com/grayson/code-clone-tool/lib/fs"
)

type performGitActionsModel struct {
	logfile      *os.File
	log          *log.Logger
	fileSystem   fs.Fs
	shouldMirror bool

	progress  progress.Model
	total     int
	completed atomic.Int32
	next      atomic.Int32

	activeActions [4]*lib.Action
	actions       []lib.Action
	api           git.GitApi
	state         performingGitActionsState
}

type doingGitActionMsg struct {
	index            int
	concurrencyIndex int
}

type finishedGitActionMsg struct {
	concurrencyIndex int
	task             lib.Task
	url              string
}

type finishedPerformingGitActions struct{}

type performingGitActionsState int

const (
	waitingToPerformGitActionsState performingGitActionsState = iota
	updatingPerformingGitActionsState
	finishedPerformingGitActionsState
)

func NewPerformGitActionsModel(fileSystem fs.Fs) *performGitActionsModel {
	file, _ := tea.LogToFile("code-clone-tool.log", "debug")
	progress := progress.New()

	return &performGitActionsModel{
		fileSystem: fileSystem,
		log:        log.Default(),
		logfile:    file,
		progress:   progress,
	}
}

func (m *performGitActionsModel) Dispose() {
	m.logfile.Close()
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
		return m, performGitAction(actual.index, actual.concurrencyIndex, m)
	case finishedGitActionMsg:
		m.completed.Add(1)
		if int(m.completed.Load()) == m.total {
			m.state = finishedPerformingGitActionsState
			return m, func() tea.Msg { return finishedPerformingGitActions{} }
		}

		next := int(m.next.Add(1))
		var cmd tea.Cmd
		if next < m.total {
			cmd = func() tea.Msg {
				return doingGitActionMsg{
					index:            next,
					concurrencyIndex: actual.concurrencyIndex,
				}
			}
		} else {
			m.activeActions[actual.concurrencyIndex] = nil
		}
		return m, cmd
	}

	return m, nil
}

func (m *performGitActionsModel) View() string {
	if m.state == waitingToPerformGitActionsState {
		return ""
	}
	completed := m.completed.Load()

	toUI := func(t lib.Task) string {
		switch t {
		case lib.Pull:
			return "pulling"
		case lib.Clone:
			return "cloning"
		}
		return ""
	}

	lines := ""
	for idx := 0; idx < 4; idx++ {
		act := m.activeActions[idx]
		if act == nil {
			continue
		}
		lines = fmt.Sprintf("%v> %v %v\n", lines, toUI(act.Task), path.Base(act.GitUrl))
	}

	return fmt.Sprintf("%v\n%v of %v finished\n%v", lines, completed, m.total, m.progress.ViewAs(float64(completed)/float64(m.total)))
}

func startGitActions(m *performGitActionsModel) tea.Cmd {
	m.next.Store(4)
	return tea.Batch(
		func() tea.Msg { return doingGitActionMsg{0, 0} },
		func() tea.Msg { return doingGitActionMsg{1, 1} },
		func() tea.Msg { return doingGitActionMsg{2, 2} },
		func() tea.Msg { return doingGitActionMsg{3, 3} },
	)
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

func performGitAction(index int, concurrencyIndex int, m *performGitActionsModel) tea.Cmd {
	actions := m.actions
	api := m.api

	return func() tea.Msg {
		var err error
		var result string
		action := actions[index]
		task := action.Task
		switch action.Task {
		case lib.Clone:
			result, err = api.Clone(action.GitUrl, action.Path)
		case lib.Pull:
			result, err = api.Pull(action.Path)
		default:
			err = fmt.Errorf("unexpected task: %v", action.Task.String())
		}
		if err != nil {
			return reportError(err)
		}
		m.log.Println(result)

		m.activeActions[concurrencyIndex] = &action
		return finishedGitActionMsg{
			concurrencyIndex: concurrencyIndex,
			task:             task,
			url:              action.GitUrl,
		}
	}
}

func determineGitClient(isMirror bool, log *log.Logger) git.GitApi {
	if isMirror {
		return git.CreateMirrorClient(log)
	}
	return git.CreateGitClient(log)
}
