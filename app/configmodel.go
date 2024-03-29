package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grayson/code-clone-tool/lib"
	"github.com/grayson/code-clone-tool/lib/fs"
)

type configmodel struct {
	textInput           textinput.Model
	shouldShowTextInput bool

	fileSystem fs.Fs

	workingDirectory    string
	personalAccessToken string
	isMirror            bool
	isComplete          bool
	url                 string
}

func NewConfigModel(env *lib.Env, fileSystem fs.Fs) *configmodel {
	ti := textinput.New()

	return &configmodel{
		textInput:           ti,
		fileSystem:          fileSystem,
		workingDirectory:    env.WorkingDirectory,
		personalAccessToken: env.PersonalAccessToken,
		url:                 env.ApiUrl,
		isMirror:            env.IsMirror.IsTruthy(),
	}
}

type cwdMsg string
type showTextInputMsg struct{}

func (c *configmodel) Init() tea.Cmd {
	return func() tea.Msg {
		if c.workingDirectory == "" {
			path, err := c.fileSystem.GetWorkingDirectory()
			if err != nil {
				return reportError(err)
			}
			return cwdMsg(path)
		}

		err := c.fileSystem.ChangeWorkingDirectory(c.workingDirectory)
		if err != nil {
			return reportError(err)
		}
		return cwdMsg(c.workingDirectory)
	}
}

func (c *configmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch actual := msg.(type) {
	case cwdMsg:
		c.workingDirectory = string(actual)
		return c, nil
	case showTextInputMsg:
		c.shouldShowTextInput = true
	case tea.KeyMsg:
		switch actual.String() {
		case tea.KeyEnter.String():
			c.shouldShowTextInput = false

			//commit
			if len(c.personalAccessToken) == 0 {
				c.personalAccessToken = c.textInput.Value()
			} else {
				c.url = c.textInput.Value()
			}
		}
	}

	if c.shouldShowTextInput {
		var cmd tea.Cmd
		c.textInput, cmd = c.textInput.Update(msg)
		return c, cmd
	}

	return c, getNextCmd(c)
}

func (c *configmodel) View() string {
	var sb strings.Builder

	hasPat, hasUrl, isMirror := " ", " ", " "
	if len(c.personalAccessToken) != 0 {
		hasPat = "✓"
	}
	if len(c.url) != 0 {
		hasUrl = "✓"
	}

	if c.isMirror {
		isMirror = "✓"
	}

	fmt.Fprintf(&sb, "Working directory: %v\n", c.workingDirectory)
	fmt.Fprintf(&sb, "[%v] Has PAT, [%v] Has Url [%v] Use `git clone --mirror`", hasPat, hasUrl, isMirror)

	if c.shouldShowTextInput {
		fmt.Fprintf(&sb, "\n%v", c.textInput.View())
	}

	return sb.String()
}

func getNextCmd(c *configmodel) tea.Cmd {
	if len(c.personalAccessToken) == 0 {
		return getPATCmd(c)
	}

	if len(c.url) == 0 {
		return getUrlCmd(c)
	}

	if !c.isComplete {
		c.isComplete = true
		return func() tea.Msg {
			return configurationCompleteMsg{
				personalAccessToken: c.personalAccessToken,
				url:                 c.url,
				isMirror:            c.isMirror,
			}
		}
	}

	return nil
}

func getPATCmd(c *configmodel) tea.Cmd {
	show := func() tea.Msg {
		resetTextInput(&c.textInput, "Personal Access Token", true)
		return showTextInputMsg{}
	}
	return tea.Sequence(show, c.textInput.Focus())
}

func getUrlCmd(c *configmodel) tea.Cmd {
	show := func() tea.Msg {
		resetTextInput(&c.textInput, "API Url", false)
		return showTextInputMsg{}
	}
	return tea.Sequence(show, c.textInput.Focus())
}

func resetTextInput(textInput *textinput.Model, placeholder string, isSecure bool) {
	textInput.Reset()
	textInput.Placeholder = placeholder
	textInput.Cursor.Blink = true

	if isSecure {
		textInput.EchoMode = textinput.EchoPassword
		textInput.EchoCharacter = '*'
	} else {
		textInput.EchoMode = textinput.EchoNormal
	}
}
