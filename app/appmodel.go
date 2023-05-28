package app

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grayson/code-clone-tool/lib"
	"github.com/grayson/code-clone-tool/lib/fs"
)

type AppModel struct {
	env        *lib.Env
	version    string
	fileSystem fs.Fs

	currentWorkingDirectory string
	err                     error
}

type cwdMsg string
type errMsg error

func InitAppModel(env *lib.Env, version string, fileSystem fs.Fs) *AppModel {
	return &AppModel{
		env: env,
	}
}

func (app *AppModel) Init() tea.Cmd {
	return func() tea.Msg {
		if app.env.WorkingDirectory == "" {
			path, err := os.Getwd()
			if err != nil {
				return errMsg(err)
			}
			return cwdMsg(path)
		}

		err := app.fileSystem.ChangeWorkingDirectory(app.env.WorkingDirectory)
		if err != nil {
			return errMsg(err)
		}
		return cwdMsg(app.env.WorkingDirectory)
	}
}

func (app *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch actual := msg.(type) {
	case tea.KeyMsg:
		return handleKeyboardEvent(actual, app)
	case cwdMsg:
		app.currentWorkingDirectory = string(actual)
		return app, nil
	case errMsg:
		app.err = error(actual)
		return app, tea.Quit
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

	fmt.Fprintf(&sb, "Working directory: %v", app.currentWorkingDirectory)

	return sb.String()
}

func handleKeyboardEvent(msg tea.KeyMsg, app *AppModel) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return app, tea.Quit
	}
	return app, nil
}
