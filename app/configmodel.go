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
	url                 string
	err                 error
}

func NewConfigModel(env *lib.Env, fileSystem fs.Fs) *configmodel {
	ti := textinput.New()

	return &configmodel{
		textInput:           ti,
		fileSystem:          fileSystem,
		workingDirectory:    env.WorkingDirectory,
		personalAccessToken: env.PersonalAccessToken,
		url:                 env.ApiUrl,
	}
}

type cwdMsg string
type patMsg string
type urlMsg string
type showTextInputMsg struct{}
type hideTextInputMsg struct{}

func (c *configmodel) Init() tea.Cmd {
	return func() tea.Msg {
		if c.workingDirectory == "" {
			path, err := c.fileSystem.GetWorkingDirectory()
			if err != nil {
				return errMsg(err)
			}
			return cwdMsg(path)
		}

		err := c.fileSystem.ChangeWorkingDirectory(c.workingDirectory)
		if err != nil {
			return errMsg(err)
		}
		return cwdMsg(c.workingDirectory)
	}
}

func (c *configmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch actual := msg.(type) {
	case cwdMsg:
		c.workingDirectory = string(actual)
		return c, nil
	case errMsg:
		c.err = error(actual)
		return c, tea.Quit
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

	hasPat, hasUrl := " ", " "
	if len(c.personalAccessToken) != 0 {
		hasPat = "✓"
	}
	if len(c.url) != 0 {
		hasUrl = "✓"
	}

	fmt.Fprintf(&sb, "Working directory: %v\n", c.workingDirectory)
	fmt.Fprintf(&sb, "[%v] Has PAT, [%v] Has Url\n", hasPat, hasUrl)

	if c.shouldShowTextInput {
		fmt.Fprintln(&sb, c.textInput.View())
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

	return nil
}

func getPATCmd(c *configmodel) tea.Cmd {
	show := func() tea.Msg {
		resetTextInput(c.textInput, "Personal Access Token", true)
		return showTextInputMsg{}
	}
	return tea.Sequence(show, c.textInput.Focus())
}

func getUrlCmd(c *configmodel) tea.Cmd {
	show := func() tea.Msg {
		resetTextInput(c.textInput, "API Url", false)
		return showTextInputMsg{}
	}
	return tea.Sequence(show, c.textInput.Focus())
}

func resetTextInput(textInput textinput.Model, placeholder string, isSecure bool) {
	textInput.Reset()
	textInput.Placeholder = placeholder
	textInput.Cursor.Blink = true

	if isSecure {
		textInput.EchoMode = textinput.EchoPassword
		textInput.EchoCharacter = '*'
	}
}
