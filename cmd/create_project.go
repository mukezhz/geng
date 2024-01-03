package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/terminal"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var newProjectCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a new project",
	Args:  cobra.MaximumNArgs(1),
	Run:   createProject,
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
