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

	infraGen.Directory = projectPath
	infras, infrasTmpl := infraGen.PathGen()

	questions := []terminal.ProjectQuestion{
		terminal.NewCheckboxQuestion(constant.InfrastructureNameKEY, "Select the infrastructure? [<space> to select]", infras),
	}

	terminal.StartInteractiveTerminal(questions)
	for _, q := range questions {
		if q.Input.Exited() {
			color.Redln("exited without completing...")
			return
		}
	}

	var selectedItems []int
	var selectedFunctions []string
	var selectedInfras []string
	for _, q := range questions {
		if q.Key == constant.InfrastructureNameKEY {
			selected := q.Input.Selected()
			for s := range selected {
				funcPath := utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "infrastructure", infrasTmpl[s]))
				funcDecl := utility.GetFunctionDeclarations(funcPath, templates.FS)
				selectedFunctions = append(selectedFunctions, funcDecl...)
				selectedItems = append(selectedItems, s)
				selectedInfras = append(selectedInfras, infras[s])
			}
		}
	}

	if err := infraGen.Validate(); err != nil {
		color.Red.Println(err)
		return
	}

	if len(selectedItems) == 0 {
		color.Red.Println("No infrastructure selected")
		return
	}

	servicesTmplMap := infraGen.Generate(data, selectedItems, selectedFunctions)
	serviceModulePath := filepath.Join(data.PackageName, "pkg", "services", "module.go")

	var servicesTmpl []string
	for k := range servicesTmplMap {
		servicesTmpl = append(servicesTmpl, k)
	}

	addService(questions, servicesTmpl, serviceModulePath, data, true, templates.FS)

	utility.PrintColorizeInfrastructureDetail(data, selectedInfras)
}
