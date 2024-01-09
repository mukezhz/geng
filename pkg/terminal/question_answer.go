package terminal

import (
	"github.com/charmbracelet/lipgloss"
)

func DefaultView(m *Model) string {
	current := m.questions[m.index]
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().
				Bold(true).
				Height(2).
				Render(current.Question),
			m.styles.InputField.Render(
				current.Input.View(),
			),
		),
	)
}
