package terminal

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return m.questions[m.index].Input.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &m.questions[m.index]
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", tea.KeyEsc.String(), tea.KeyEscape.String():
			current.Input.SetExited()
			return m, tea.Quit

		case "enter":
			if m.index == len(m.questions)-1 {
				m.done = true
			}
			current.Answer = current.Input.Value()
			m.Next()
			if m.done {
				return m, tea.Quit
			}
			return m, current.Input.Blur
		}
	}
	current.Input, cmd = current.Input.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.done {
		var output string
		for _, q := range m.questions {
			output += fmt.Sprintf("%s: %s\n", q.Question, q.Answer)
		}
		return output
	}
	if m.width == 0 {
		return "loading..."
	}
	return DefaultView(m)
}

func (m *Model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func StartInteractiveTerminal(questions []ProjectQuestion) {
	m := New(questions)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
