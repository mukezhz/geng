package cmd

import (
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/gen"
	"github.com/mukezhz/geng/pkg/terminal"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service add [name]",
	Short: "Add a new Service",
	Args:  cobra.MaximumNArgs(2),
	Run:   addServiceHandler,
}

var serviceGen = gen.ServiceGenerator{}

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

	serviceGen.Directory = projectPath
	choice := serviceGen.GetChoices()

	questions := []terminal.ProjectQuestion{
		terminal.NewCheckboxQuestion(constant.ServiceNameKEY, "Select the Service? [<space> to select]", choice.Items),
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
		if q.Key == constant.ServiceNameKEY {
			selected := q.Input.Selected()
			for s := range selected {
				selectedItems = append(selectedItems, s)
			}
		}
	}

	if err := serviceGen.Generate(data, selectedItems); err != nil {
		color.Red.Printf("Generation error: %v\n", err)
		return
	}

	selectedServices := serviceGen.GetSelectedItems(selectedItems)
	utility.PrintColorizeServiceDetail(data, selectedServices)
}
