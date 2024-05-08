package gen

import (
	"errors"
	"fmt"
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
	templates []string
	infra     []string
}

// PathGen generates necessary path argument for infrastructure generation
// genration is carried out with directory argument
func (i *InfraGenerator) PathGen() ([]string, []string) {
	i.modPath = filepath.Join(i.Directory, "pkg", "infrastructure", "module.go")
	i.infraPath = utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "infrastructure"))
	i.templates = utility.ListDirectory(templates.FS, i.infraPath)

	replaceFunc := func(q string) string {
		return strings.Replace(q, ".tmpl", "", 1)
	}

	i.infra = utility.Map[string, string](i.templates, replaceFunc)

	return i.infra, i.templates
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
	items []int, // selected items for generation
	functions []string, // selected functions for generation
) map[string]bool {

	g.PathGen()

	if len(items) == 0 {
		return nil
	}

  fmt.Println(functions)

	updatedCode := utility.AddListOfProvideInFxOptions(g.modPath, functions)
	utility.WriteContentToPath(g.modPath, updatedCode)

	servicesTmplMap := make(map[string]bool)
	for _, item := range items {
		currTemplate := g.templates[item]
		templatePath := filepath.Join(".", "templates", "wesionary", "infrastructure", currTemplate)
		templatePath = utility.IgnoreWindowsPath(templatePath)

		targetRoot := filepath.Join(data.Directory, "pkg", "infrastructure", strings.Replace(currTemplate, ".tmpl", ".go", 1))

		fileName := strings.Replace(currTemplate, ".tmpl", "", 1)
		serviceTemplatePath := utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "service"))
		for _, file := range utility.ListDirectory(templates.FS, serviceTemplatePath) {
			if strings.Contains(file, fileName) {
				servicesTmplMap[file] = true
			}
		}
		utility.GenerateFromEmbeddedTemplate(templates.FS, templatePath, targetRoot, data)
	}

	return servicesTmplMap
}
