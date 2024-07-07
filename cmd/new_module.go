package cmd

import (
	"path/filepath"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/model"
	"github.com/mukezhz/geng/pkg/terminal"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/mukezhz/geng/templates"
	"github.com/spf13/cobra"
)

var newModuleCmd = &cobra.Command{
	Use:   "gen mod [name]",
	Short: "Create a new domain",
	Long: `
Create a new module|service|middleware in the project.
Example: 
  geng gen mod [name]
  geng gen srv [name]
  geng gen mid [name]

Default:
  geng gen -> geng gen mod
	`,
	Args: cobra.MaximumNArgs(2),
	Run:  generate,
}

func generate(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		args = append(args, "module")
	}

	projectPath, projectModule := getProjectPath()
	if projectModule == nil || projectPath == "" {
		return
	}
	mainModulePath := filepath.Join(projectPath, "domain", "module.go")
	if !utility.FileExists(mainModulePath) {
		color.Redln("Error: module.go not found")
		return
	}
	// Define the directory structure
	generateModule(mainModulePath, args, *projectModule)

}

func generateModule(mainModulePath string, args []string, projectModule model.GoMod) {
	var moduleName string
	if len(args) == 1 {
		questions := []terminal.ProjectQuestion{
			terminal.NewShortQuestion(constant.ModueleNameKEY, constant.ModueleNameKEY+" *", "Enter Module Name:"),
		}
		terminal.StartInteractiveTerminal(questions)

		for _, q := range questions {
			switch q.Key {
			case constant.ModueleNameKEY:
				moduleName = q.Answer
			}
			if q.Input.Exited() {
				color.Redln("exited without completing...")

			}
		}
	} else {
		moduleName = args[1]
	}
	if !utility.CheckGolangIdentifier(moduleName) {
		color.Redln("Error: module name is invalid")
		return
	}
	data := utility.GetModuleDataFromModuleName(moduleName, projectModule.Module, projectModule.GoVersion)
	data.IsModuleGenerated = true

	targetRoot := filepath.Join(".", "domain", data.PackageName)
	templatePath := utility.IgnoreWindowsPath(filepath.Join(".", "templates", "wesionary", "module"))

	err := utility.GenerateFiles(templates.FS, templatePath, targetRoot, data)
	if err != nil {
		color.Redln("Error: generate file", err)
		return
	}

	updatedCode := utility.AddAnotherFxOptionsInModule(mainModulePath, data.PackageName, data.ProjectModuleName)
	utility.WriteContentToPath(mainModulePath, updatedCode)

	utility.PrintColorizeModuleDetail(data)

}
