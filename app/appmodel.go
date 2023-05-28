package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grayson/code-clone-tool/lib"
	"github.com/grayson/code-clone-tool/lib/fs"
)

type AppModel struct {
	env     *lib.Env
	version string

	err      error
	children []tea.Model
}

type errMsg error
type configurationCompleteMsg struct{}

func InitAppModel(env *lib.Env, version string, fileSystem fs.Fs) *AppModel {
	return &AppModel{
		env: env,
		children: []tea.Model{
			NewConfigModel(env, fileSystem),
		},
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
	}

	for _, model := range app.children {
		_, cmd := model.Update(msg)
		if cmd != nil {
			return app, cmd
		}
	}

	return app, nil
}

func (app *AppModel) View() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "code-clone-tool %v\n", app.version)
	if app.err != nil {
		fmt.Fprintf(&sb, "%v", app.err)
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
