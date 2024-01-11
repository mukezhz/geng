package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/terminal"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var newModuleCmd = &cobra.Command{
	Use:   "gen module [name]",
	Short: "Create a new domain",
	Args:  cobra.MaximumNArgs(2),
	Run:   createModule,
}

func createModule(_ *cobra.Command, args []string) {
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
		}
	} else {
		moduleName = args[1]
	}
	if !utility.CheckGolangIdentifier(moduleName) {
		color.Redln("Error: module name is invalid")
		return
	}
	data := utility.GetModuleDataFromModuleName(moduleName, projectModule.Module, projectModule.GoVersion)

	// Define the directory structure
	targetRoot := filepath.Join(".", "domain", data.PackageName)
	templatePath := filepath.Join(".", "templates", "wesionary", "module")

	err = utility.GenerateFiles(templatesFS, templatePath, targetRoot, data)
	if err != nil {
		color.Redln("Error: generate file", err)
		return
	}

	updatedCode := utility.AddAnotherFxOptionsInModule(mainModulePath, data.PackageName, data.ProjectModuleName)
	utility.WriteContentToPath(mainModulePath, updatedCode)

	utility.PrintColorizeModuleDetail(data)
}
