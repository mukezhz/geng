package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/gen"
	"github.com/mukezhz/geng/pkg/terminal"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/mukezhz/geng/templates"
	"github.com/spf13/cobra"
)

var infraCmd = &cobra.Command{
	Use:   "infra add [name]",
	Short: "Add a new infrastructure",
	Args:  cobra.MaximumNArgs(2),
	Run:   addInfrastructureHandler,
}

var infraGen = gen.InfraGenerator{}

func addInfrastructureHandler(_ *cobra.Command, args []string) {

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

	// generate in project path
	infraGen.Directory = projectPath

	choice := infraGen.GetChoices()

	questions := []terminal.ProjectQuestion{
		terminal.NewCheckboxQuestion(constant.InfrastructureNameKEY, "Select the infrastructure? [<space> to select]", choice.Items),
	}

	terminal.StartInteractiveTerminal(questions)
	for _, q := range questions {
		if q.Input.Exited() {
			color.Redln("exited without completing...")
			return
		}
	}

	// get selected items from questions
	var selectedItems []int
	for _, q := range questions {
		if q.Key == constant.InfrastructureNameKEY {
			selected := q.Input.Selected()
			for s := range selected {
				selectedItems = append(selectedItems, s)
			}
		}
	}

	if err := infraGen.Validate(); err != nil {
		color.Red.Println(err)
		return
	}

	if err := infraGen.Generate(data, selectedItems); err != nil {
		color.Red.Printf("Generation error: %v\n", err)
		return
	}

  // TODO: refactor service generation logic
	var servicesTmpl []string
	for _, item := range selectedItems {
		currTemplate := choice.Templates[item]
		fileName := strings.Replace(currTemplate, ".tmpl", "", 1)
		serviceTemplatePath := utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "service"))
		for _, file := range utility.ListDirectory(templates.FS, serviceTemplatePath) {
			// weird logic for now
			if strings.Contains(file, fileName) {
				servicesTmpl = append(servicesTmpl, file)
			}
		}
	}

	serviceModulePath := filepath.Join(data.PackageName, "pkg", "services", "module.go")
	addService(questions, servicesTmpl, serviceModulePath, data, true, templates.FS)

	selected := infraGen.GetSelectedItems(selectedItems)
	utility.PrintColorizeInfrastructureDetail(data, selected)
}
