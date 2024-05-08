package cmd

import (
	"path/filepath"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/gen"
	"github.com/mukezhz/geng/pkg/terminal"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/mukezhz/geng/templates"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a new project",
	Args:  cobra.MaximumNArgs(1),
	Run:   createProject,
}

var projectGen = gen.ProjectGenerator{}

func init() {
	projectCmd.Flags().StringVarP(&projectGen.ModuleName, "mod", "m", "", "module name")
	projectCmd.Flags().StringVarP(&projectGen.Directory, "dir", "d", "", "target directory")
	projectCmd.Flags().StringVarP(&projectGen.GoVersion, "version", "v", "", "version support: Default: 1.20")
}

func createProject(cmd *cobra.Command, args []string) {
	var questions []terminal.ProjectQuestion

	infraGen := gen.InfraGenerator{Directory: "."}
	infras, infrasTmpl := infraGen.PathGen()

	if len(args) == 0 {
		questions = []terminal.ProjectQuestion{
			terminal.NewShortQuestion(constant.ProjectNameKEY, constant.ProjectName+" *", "Enter Project Name:"),
			terminal.NewShortQuestion(constant.ProjectModuleNameKEY, constant.ProjectModuleName+" *", "Enter Module Name:"),
			terminal.NewShortQuestion(constant.AuthorKEY, constant.Author+" [Optional]", "Enter Author Detail[Mukesh Chaudhary <mukezhz@duck.com>] [Optional]"),
			terminal.NewLongQuestion(constant.ProjectDescriptionKEY, constant.ProjectDescription+" [Optional]", "Enter Project Description [Optional]"),
			terminal.NewShortQuestion(constant.GoVersionKEY, constant.GoVersion+" [Optional]", "Enter Go Version (Default: 1.20) [Optional]"),
			terminal.NewShortQuestion(constant.DirectoryKEY, constant.Directory+" [Optional]", "Enter Project Directory (Default: package_name) [Optional]"),
			terminal.NewCheckboxQuestion(constant.InfrastructureNameKEY, "Select the infrastructure? [<space> to select] [Optional]", infras),
		}

		terminal.StartInteractiveTerminal(questions)

		questionMap := make(map[string]string)

		for _, q := range questions {
			questionMap[q.Key] = q.Answer
			if q.Input.Exited() {
				color.Redln("exited without completing...")
				return
			}
		}

		projectGen.Fill(questionMap)

	} else {
		projectGen.Name = args[0]
		projectGen.ModuleName, _ = cmd.Flags().GetString("mod")
		projectGen.GoVersion, _ = cmd.Flags().GetString("version")
		projectGen.Directory, _ = cmd.Flags().GetString("dir")
	}

	if err := projectGen.Validate(); err != nil {
		color.Redln(err.Error())
		return
	}

	data, err := projectGen.Generate()
	if err != nil {
		color.Redln(err.Error())
		return
	}

	// TODO: refactor, selected items generation logic from question generation
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

  infraGen.Directory = data.Directory
	if err := infraGen.Validate(); err != nil {
		color.Red.Println(err)
		return
	}

	if len(selectedItems) == 0 {
		color.Red.Println("No infrastructure selected")
		return
	}

	servicesTmplMap := infraGen.Generate(*data, selectedItems, selectedFunctions)
	serviceModulePath := filepath.Join(data.PackageName, "pkg", "services", "module.go")

	var servicesTmpl []string
	for k := range servicesTmplMap {
		servicesTmpl = append(servicesTmpl, k)
	}

	addService(questions, servicesTmpl, serviceModulePath, *data, true, templates.FS)

	utility.PrintColorizeInfrastructureDetail(*data, selectedInfras)

}
