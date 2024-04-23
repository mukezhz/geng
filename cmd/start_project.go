package cmd

import (
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var startProjectCmd = &cobra.Command{
	Use:   "start [project name]",
	Short: "Execute the project",
	Long: `
Execute the project:
Alias to "go run main.go app:serve"

For available command: see "geng project"
	`,
	Args: cobra.MaximumNArgs(1),
	Run:  startProject,
}

func startProject(_ *cobra.Command, args []string) {
	program := "go"
	if len(args) == 0 {
		args = append(args, "app:serve")
	}
	commands := []string{"run", "main.go"}
	// execute command from golang
	err := utility.ExecuteCommand(program, commands, args...)
	if err != nil {
		return
	}
}
