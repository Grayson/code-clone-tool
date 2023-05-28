package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grayson/code-clone-tool/lib"
	"github.com/grayson/code-clone-tool/lib/fs"
)

type configmodel struct {
	fileSystem fs.Fs

	workingDirectory    string
	personalAccessToken string
	url                 string
	err                 error
}

func NewConfigModel(env *lib.Env, fileSystem fs.Fs) *configmodel {
	return &configmodel{
		fileSystem:          fileSystem,
		workingDirectory:    env.WorkingDirectory,
		personalAccessToken: env.PersonalAccessToken,
		url:                 env.ApiUrl,
	}
}

type cwdMsg string
type patMsg string
type urlMsg string

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
	}
	return c, nil
}

func (c *configmodel) View() string {
	return fmt.Sprintf("Working directory: %v", c.workingDirectory)
}
