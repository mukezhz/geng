package terminal

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type shortAnswerField struct {
	textinput textinput.Model
}

func NewShortAnswerField(p string) *shortAnswerField {
	a := shortAnswerField{}
	if p == "" {
		p = "Your Answer here"
	}
	model := textinput.New()
	model.Placeholder = p
	model.Focus()

	a.textinput = model
	return &a
}

func (a *shortAnswerField) Blink() tea.Msg {
	return textinput.Blink()
}

func (a *shortAnswerField) Init() tea.Cmd {
	return nil
}

func (a *shortAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textinput, cmd = a.textinput.Update(msg)
	return a, cmd
}

func (a *shortAnswerField) View() string {
	return a.textinput.View()
}

func (a *shortAnswerField) Focus() tea.Cmd {
	return a.textinput.Focus()
}

func (a *shortAnswerField) SetValue(s string) {
	a.textinput.SetValue(s)
}

func (a *shortAnswerField) Blur() tea.Msg {
	return a.textinput.Blur
}

func (a *shortAnswerField) Value() string {
	return a.textinput.Value()
}
