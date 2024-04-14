package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/model"
	"github.com/mukezhz/geng/pkg/terminal"
	"github.com/mukezhz/geng/pkg/utility"
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
	projectModule, err := utility.GetModuleNameFromGoModFile()
	if err != nil {
		fmt.Println("Error finding Module name from go.mod:", err)
		return
	}
	currentDir, err := os.Getwd()
	if err != nil {
		color.Redln("Error getting current directory:", err)
		panic(err)
	}
	projectPath, err := utility.FindGitRoot(currentDir)
	if err != nil {
		fmt.Println("Error finding Git root:", err)
		return
	}
	// Define the directory structure
	generateModule(projectPath, args, projectModule)

}

func generateModule(projectPath string, args []string, projectModule model.GoMod) {
	mainModulePath := filepath.Join(projectPath, "domain", "module.go")
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

	targetRoot := filepath.Join(".", "domain", data.PackageName)
	templatePath := filepath.Join(".", "templates", "wesionary", "module")

	err := utility.GenerateFiles(templatesFS, templatePath, targetRoot, data)
	if err != nil {
		color.Redln("Error: generate file", err)
		return
	}

	updatedCode := utility.AddAnotherFxOptionsInModule(mainModulePath, data.PackageName, data.ProjectModuleName)
	utility.WriteContentToPath(mainModulePath, updatedCode)

	utility.PrintColorizeModuleDetail(data)

}
