package gen

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/mukezhz/geng/pkg/model"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/mukezhz/geng/templates"
)

type ServiceGenerator struct {
	Directory string

	modPath   string
	infraPath string

	choice *ServiceChoice
}

type ServiceChoice struct {
	Items     []string
	Templates []string
}

// GetSelectedItems converts selected items (integer) from template into
// array of strings depending on the choices generated.
func (i *ServiceGenerator) GetSelectedItems(selectedItems []int) []string {
	retItem := make([]string, len(selectedItems))
	items := i.GetChoices().Items
	for i, selectedIndex := range selectedItems {
		retItem[i] = items[selectedIndex]
	}

	return retItem
}

func (s *ServiceGenerator) GetChoices() *ServiceChoice {
	s.modPath = filepath.Join(s.Directory, "pkg", "services", "module.go")
	s.infraPath = utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "service"))
	templates := utility.ListDirectory(templates.FS, s.infraPath)

	replaceFunc := func(q string) string {
		return strings.Replace(q, ".tmpl", "", 1)
	}

	services := utility.Map[string, string](templates, replaceFunc)

	s.choice = &ServiceChoice{
		Items:     services,
		Templates: templates,
	}

	return s.choice
}

// SimilarChoice gives similar choice to passed values from templates.
// simple strings.Contain operation is carried out.
func (s *ServiceGenerator) SimilarChoice(shouldMatch []string) []int {
	var matches []int
	choices := s.GetChoices()

	for _, m := range shouldMatch {
		for i, item := range choices.Items {
      // similarity check for now
			if strings.Contains(item, m) {
				matches = append(matches, i)
			}
		}
	}

	return matches
}

func (s *ServiceGenerator) Generate(
	data model.ModuleData,
	selectedItems []int,
) error {

	s.GetChoices()

	if len(s.choice.Items) == 0 {
		return errors.New("no choice to choose from, templates mismatch")
	}

	if len(selectedItems) == 0 {
		return errors.New("noting choosen by the user, skipping")
	}

	var functions []string
	for _, index := range selectedItems {
		funcPath := utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "service", s.choice.Templates[index]))
		funcs := utility.GetFunctionDeclarations(funcPath, templates.FS)

		filterFunc := func(q string) bool {
			return strings.Contains(q, "New")
		}
		filteredServices := utility.Filter[string](funcs, filterFunc)

		functions = append(functions, filteredServices...)
	}

	updatedCode := utility.AddListOfProvideInFxOptions(s.modPath, functions)
	utility.WriteContentToPath(s.modPath, updatedCode)

	for _, item := range selectedItems {
		currTemplate := s.choice.Templates[item]
		templatePath := filepath.Join(".", "templates", "wesionary", "service", currTemplate)
		templatePath = utility.IgnoreWindowsPath(templatePath)

		targetRoot := filepath.Join(data.Directory, "pkg", "services", strings.Replace(currTemplate, ".tmpl", ".go", 1))

		utility.GenerateFromEmbeddedTemplate(templates.FS, templatePath, targetRoot, data)
	}

	return nil
}
