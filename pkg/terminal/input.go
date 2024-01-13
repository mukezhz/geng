package terminal

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Blink() tea.Msg
	Blur() tea.Msg
	Focus() tea.Cmd
	SetValue(string)
	Value() string
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
	Selected() map[int]any
	SetExited()
	Exited() bool
}

// We need to have a wrapper for our bubbles as they don't currently implement the tea.Model interface
// textinput, textarea
