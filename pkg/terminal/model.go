package terminal

type Model struct {
	styles    *Styles
	index     int
	questions []ProjectQuestion
	width     int
	height    int
	done      bool
}

type ProjectQuestion struct {
	Key         string
	PlaceHolder string
	Question    string
	Answer      string
	Input       Input
}

func NewQuestion(k, q, p string) ProjectQuestion {
	return ProjectQuestion{Key: k, PlaceHolder: p, Question: q}
}

func NewShortQuestion(k, q, p string) ProjectQuestion {
	question := NewQuestion(k, q, p)
	model := NewShortAnswerField(p)
	question.Input = model
	return question
}

func NewLongQuestion(k, q, p string) ProjectQuestion {
	question := NewQuestion(k, q, p)
	model := NewLongAnswerField()
	question.Input = model
	return question
}

func NewCheckboxQuestion(k, t string, items []string) ProjectQuestion {
	question := NewQuestion(k, t, "")
	model := NewCheckBoxField(t, items)
	question.Input = model
	return question
}

func New(questions []ProjectQuestion) *Model {
	styles := DefaultStyles()
	return &Model{styles: styles, questions: questions}
}
