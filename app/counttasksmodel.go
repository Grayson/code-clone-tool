package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grayson/code-clone-tool/lib"
)

type counttasksmodel struct {
	pulls  int
	clones int

	isVisible bool
}

func (c *counttasksmodel) Init() tea.Cmd {
	return nil
}

func (c *counttasksmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch actual := msg.(type) {
	case repoResponseMsg:
		c.isVisible = true
	case finishedGitActionMsg:
		switch actual.task {
		case lib.Clone:
			c.clones++
		case lib.Pull:
			c.pulls++
		}
	}
	return c, nil
}

func (c *counttasksmodel) View() string {
	if !c.isVisible {
		return ""
	}
	return fmt.Sprintf("%v Pulls; %v Clones", c.pulls, c.clones)
}
