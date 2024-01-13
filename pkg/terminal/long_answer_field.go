package terminal

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type longAnswerField struct {
	textarea textarea.Model
	exited   bool
}

func NewLongAnswerField() *longAnswerField {
	a := longAnswerField{}

	model := textarea.New()
	model.Placeholder = "Your Answer here"
	model.Focus()

	a.textarea = model
	return &a
}

func (a *longAnswerField) Blink() tea.Msg {
	return textarea.Blink()
}

func (a *longAnswerField) Init() tea.Cmd {
	return nil
}

func (a *longAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textarea, cmd = a.textarea.Update(msg)
	return a, cmd
}

func (a *longAnswerField) View() string {
	return a.textarea.View()
}

func (a *longAnswerField) Focus() tea.Cmd {
	return a.textarea.Focus()
}

func (a *longAnswerField) SetValue(s string) {
	a.textarea.SetValue(s)
}

func (a *longAnswerField) Blur() tea.Msg {
	return a.textarea.Blur
}

func (a *longAnswerField) Value() string {
	return a.textarea.Value()
}

func (a *longAnswerField) Selected() map[int]any {
	return nil
}

func (a *longAnswerField) SetExited() {
	a.exited = true
}

func (a *longAnswerField) Exited() bool {
	return a.exited
}
