package cmd

import (
	"embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/model"
	"github.com/mukezhz/geng/pkg/terminal"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var addServiceCmd = &cobra.Command{
	Use:   "service add [name]",
	Short: "Add a new Service",
	Args:  cobra.MaximumNArgs(2),
	Run:   addServiceHandler,
}

func addServiceHandler(_ *cobra.Command, args []string) {
	if len(args) > 0 && !strings.Contains(args[0], "add") {
		color.Redln("Error: invalid command")
		return
	}
	projectModule, err := utility.GetModuleNameFromGoModFile()
	if err != nil {
		color.Redln("Error finding Module name from go.mod:", err)
		return
	}
	data := utility.GetModuleDataFromModuleName("", projectModule.Module, projectModule.GoVersion)
	currentDir, err := os.Getwd()
	if err != nil {
		color.Redln("Error getting current directory:", err)
		panic(err)
	}
	projectPath, err := utility.FindGitRoot(currentDir)
	if err != nil {
		color.Redln("Error finding Git root:", err)
		return
	}
	serviceModulePath := filepath.Join(projectPath, "pkg", "services", "module.go")
	templateInfraPath := utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "service"))
	servicesTmpl := utility.ListDirectory(templatesFS, templateInfraPath)
	services := utility.Map[string, string](servicesTmpl, func(q string) string {
		return strings.Replace(q, ".tmpl", "", 1)
	})

	questions := []terminal.ProjectQuestion{
		terminal.NewCheckboxQuestion(constant.ServiceNameKEY, "Select the Service? [<space> to select]", services),
	}

	terminal.StartInteractiveTerminal(questions)
	for _, q := range questions {
		if q.Input.Exited() {
			color.Redln("exited without completing...")
			return
		}
	}

	items := addService(questions, servicesTmpl, serviceModulePath, data, false, templatesFS)
	if len(items) == 0 {
		color.Red.Println("No Service selected")
		return
	}
	var selectedServices []string
	for _, i := range items {
		selectedServices = append(selectedServices, services[i])
	}
	utility.PrintColorizeServiceDetail(data, selectedServices)
}

func addService(
	questions []terminal.ProjectQuestion,
	servicesTmpl []string,
	serviceModulePath string,
	data model.ModuleData,
	isNewProject bool,
	templatesFS embed.FS,
) []int {
	var functions []string
	var items []int
	for _, q := range questions {
		switch q.Key {
		case constant.ServiceNameKEY:
			selected := q.Input.Selected()
			for s := range selected {
				funcs := utility.GetFunctionDeclarations(utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "service", servicesTmpl[s])), templatesFS)
				filteredServices := utility.Filter[string](funcs, func(q string) bool {
					return strings.Contains(q, "New")
				})
				functions = append(functions,
					filteredServices...,
				)
				items = append(items, s)
			}
		case constant.InfrastructureNameKEY:
			for i, s := range servicesTmpl {
				funcs := utility.GetFunctionDeclarations(utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "service", s)), templatesFS)
				filteredServices := utility.Filter[string](funcs, func(q string) bool {
					return strings.Contains(q, "New")
				})
				functions = append(functions,
					filteredServices...,
				)
				items = append(items, i)
			}

		}
	}
	if len(items) == 0 {
		return items
	}
	updatedCode := utility.AddListOfProvideInFxOptions(serviceModulePath, functions)
	utility.WriteContentToPath(serviceModulePath, updatedCode)

	for _, i := range items {
		templatePath := utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "service", servicesTmpl[i]))
		var targetRoot string
		if isNewProject {
			targetRoot = filepath.Join(data.PackageName, "pkg", "services", strings.Replace(servicesTmpl[i], ".tmpl", ".go", 1))
		} else {
			targetRoot = filepath.Join(".", "pkg", "services", strings.Replace(servicesTmpl[i], ".tmpl", ".go", 1))
		}
		utility.GenerateFromEmbeddedTemplate(templatesFS, templatePath, targetRoot, data)
	}
	return items
}
