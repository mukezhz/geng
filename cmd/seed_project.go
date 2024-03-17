package cmd

import (
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var seedProjectCmd = &cobra.Command{
	Use:   "seed [project name]",
	Short: "Seed the project",
	Args:  cobra.MaximumNArgs(1),
	Run:   seedProject,
}

func seedProject(_ *cobra.Command, args []string) {
	program := "go"
	commands := []string{"run", "main.go", "seed"}
	// execute command from golang
	err := utility.ExecuteCommand(program, commands, args...)
	if err != nil {
		return
	}
}
