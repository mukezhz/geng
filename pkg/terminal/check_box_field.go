package terminal

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CheckBoxField struct {
	selected map[int]any
	title    string
	choices  []string
	cursor   int
	exited   bool
}

func NewCheckBoxField(t string, choices []string) *CheckBoxField {
	m := CheckBoxField{
		choices:  choices,
		selected: make(map[int]any),
		title:    t,
	}
	return &m
}

func (m *CheckBoxField) Blink() tea.Msg {
	return nil
}

func (m *CheckBoxField) Init() tea.Cmd {
	return nil
}

func (m *CheckBoxField) Update(msg tea.Msg) (Input, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m *CheckBoxField) View() string {
	s := ""

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		style := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingTop(2).
			PaddingLeft(4).
			Width(22)
		style.Render(m.title)

		s += fmt.Sprintf("\n\n %s [%s] %s \n", cursor, checked, choice)
	}

	// Send the UI for rendering
	return s
}

func (m *CheckBoxField) Focus() tea.Cmd {
	return nil
}

func (a *CheckBoxField) SetValue(s string) {
}

func (a *CheckBoxField) Blur() tea.Msg {
	return nil
}

func (a *CheckBoxField) Value() string {
	return ""
}

func (a *CheckBoxField) Selected() map[int]any {
	return a.selected
}

func (a *CheckBoxField) SetExited() {
	a.exited = true
}

func (a *CheckBoxField) Exited() bool {
	return a.exited
}
