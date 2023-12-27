package main

import (
	"embed"
	"fmt"
	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/terminal"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

//go:embed templates/wesionary/*
var templatesFS embed.FS

var rootCmd = &cobra.Command{
	Use:              "geng",
	Short:            "Go Generate[geng] is a tool for generating Go modules",
	Long:             "Go Generate[geng] is a tool for generating Go modules. ",
	TraverseChildren: true,
}

var newModuleCmd = &cobra.Command{
	Use:   "gen features [name]",
	Short: "Create a new features",
	Args:  cobra.MaximumNArgs(2),
	Run:   createModule,
}

var newProjectCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a new project",
	Args:  cobra.MaximumNArgs(1),
	Run:   createProject,
}

var runProjectCmd = &cobra.Command{
	Use:   "run [project name]",
	Short: "Run the project",
	Args:  cobra.MaximumNArgs(1),
	Run:   runProject,
}

func init() {
	newProjectCmd.Flags().StringP("mod", "m", "", "features name")
	newProjectCmd.Flags().StringP("dir", "d", "", "target directory")
	newProjectCmd.Flags().StringP("version", "v", "", "version support")
	rootCmd.AddCommand(newModuleCmd, newProjectCmd, runProjectCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
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
	mainModulePath := filepath.Join(projectPath, "domain", "features", "module.go")
	if err != nil {
		panic(err)
	}
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
	targetRoot := filepath.Join(".", "domain", "features", data.PackageName)
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

func createProject(cmd *cobra.Command, args []string) {
	var projectName string
	var projectModuleName string
	var goVersion string
	var projectDescription string
	var author string
	var directory string

	if len(args) == 0 {
		questions := []terminal.ProjectQuestion{
			terminal.NewShortQuestion(constant.ProjectNameKEY, constant.ProjectName+" *", "Enter Project Name:"),
			terminal.NewShortQuestion(constant.ProjectModuleNameKEY, constant.ProjectModuleName+" *", "Enter Module Name:"),
			terminal.NewShortQuestion(constant.AuthorKEY, constant.Author+" [Optional]", "Enter Author Detail[Mukesh Chaudhary <mukezhz@duck.com>] [Optional]"),
			terminal.NewLongQuestion(constant.ProjectDescriptionKEY, constant.ProjectDescription+" [Optional]", "Enter Project Description [Optional]"),
			terminal.NewShortQuestion(constant.GoVersionKEY, constant.GoVersion+" [Optional]", "Enter Go Version (Default: 1.20) [Optional]"),
			terminal.NewShortQuestion(constant.DirectoryKEY, constant.Directory+" [Optional]", "Enter Project Directory (Default: package_name) [Optional]"),
		}
		terminal.StartInteractiveTerminal(questions)

		for _, q := range questions {
			switch q.Key {
			case constant.ProjectNameKEY:
				projectName = q.Answer
			case constant.ProjectDescriptionKEY:
				projectDescription = q.Answer
			case constant.AuthorKEY:
				author = q.Answer
			case constant.ProjectModuleNameKEY:
				projectModuleName = q.Answer
			case constant.GoVersionKEY:
				goVersion = q.Answer
			case constant.DirectoryKEY:
				directory = q.Answer
			}
		}
	} else {
		projectName = args[0]
		projectModuleName, _ = cmd.Flags().GetString("mod")
		goVersion, _ = cmd.Flags().GetString("version")
		directory, _ = cmd.Flags().GetString("dir")
	}

	goVersion = utility.CheckVersion(goVersion)
	if projectName == "" {
		color.Redln("Error: project name is required")
		return
	}
	if projectModuleName == "" {
		color.Redln("Error: module name is required")
		return
	}

	data := utility.GetModuleDataFromModuleName(projectName, projectModuleName, goVersion)
	data.ProjectDescription = projectDescription
	data.Author = author

	data.Directory = directory
	if data.Directory == "" {
		data.Directory = filepath.Join(data.Directory, data.PackageName)
	}
	targetRoot := data.Directory

	templatePath := filepath.Join("templates", "wesionary", "project")
	err := utility.GenerateFiles(templatesFS, templatePath, targetRoot, data)
	if err != nil {
		color.Redln("Error generate file", err)
		return
	}

	utility.PrintColorizeProjectDetail(data)
	fmt.Println("")
}

func runProject(cmd *cobra.Command, args []string) {
	runGo := "go"
	// execute command from golang
	err := utility.ExecuteCommand(runGo, args...)
	if err != nil {
		return
	}
}
