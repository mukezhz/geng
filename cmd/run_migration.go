package cmd

import (
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var migrationProjectCmd = &cobra.Command{
	Use:   "migrate [project name]",
	Short: "migrate the project",
	Args:  cobra.MaximumNArgs(0),
	Run:   migrationProject,
}

func migrationProject(_ *cobra.Command, args []string) {
	program := "go"
	commands := []string{"run", "main.go", "migrate:run"}
	// execute command from golang
	err := utility.ExecuteCommand(program, commands, args...)
	if err != nil {
		return
	}
}
