package gen

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/mukezhz/geng/pkg/model"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/mukezhz/geng/templates"
)

type InfraGenerator struct {
	Directory string

	modPath   string
	infraPath string

	choice *InfraChoice
}

type InfraChoice struct {
	Items     []string
	Templates []string
}

// GetChoices generates necessary choices for infrastructure selection
// the generated choices are filled up in the generator struct automatically.
// returns the available choices as per the templates available.
func (i *InfraGenerator) GetChoices() *InfraChoice {
	i.modPath = filepath.Join(i.Directory, "pkg", "infrastructure", "module.go")
	i.infraPath = utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "infrastructure"))

	templates := utility.ListDirectory(templates.FS, i.infraPath)

	replaceFunc := func(q string) string {
		return strings.Replace(q, ".tmpl", "", 1)
	}

	items := utility.Map[string, string](templates, replaceFunc)
	i.choice = &InfraChoice{
		Items:     items,
		Templates: templates,
	}

	return i.choice
}

// GetSelectedItems converts selected items (integer) from template into
// array of strings depending on the choices generated.
func (i *InfraGenerator) GetSelectedItems(selectedItems []int) []string {
	retItem := make([]string, len(selectedItems))
	items := i.GetChoices().Items
	for i, selectedIndex := range selectedItems {
		retItem[i] = items[selectedIndex]
	}

	return retItem
}

func (i *InfraGenerator) Validate() error {
	if i.Directory == "" {
		return errors.New("InfraGenerator: directory argument is empty")
	}

	return nil
}

// TODO: remove dependencies on model, items and functions
func (g *InfraGenerator) Generate(
	data model.ModuleData,
	selectedItems []int, // selected items for generation
) error {

	g.GetChoices()

	if len(g.choice.Items) == 0 {
		return errors.New("no choice to choose from, templates mismatch")
	}

	// nothing choosen by the user
	if len(selectedItems) == 0 {
		return errors.New("nothing choosen by the user, skipping")
	}

	// generate function declarations of selected infrastructures
	var functions []string
	for _, index := range selectedItems {
		funcPath := utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "infrastructure", g.choice.Templates[index]))
		funcDecl := utility.GetFunctionDeclarations(funcPath, templates.FS)
		functions = append(functions, funcDecl...)
	}

	updatedCode := utility.AddListOfProvideInFxOptions(g.modPath, functions)
	utility.WriteContentToPath(g.modPath, updatedCode)

	// servicesTmplMap := make(map[string]bool)
	for _, item := range selectedItems {
		currTemplate := g.choice.Templates[item]
		templatePath := filepath.Join(".", "templates", "wesionary", "infrastructure", currTemplate)
		templatePath = utility.IgnoreWindowsPath(templatePath)

		targetRoot := filepath.Join(data.Directory, "pkg", "infrastructure", strings.Replace(currTemplate, ".tmpl", ".go", 1))

		utility.GenerateFromEmbeddedTemplate(templates.FS, templatePath, targetRoot, data)
	}

	return nil
}
