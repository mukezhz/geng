package cmd

import (
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var seedProjectCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the project",
	Run:   seedProject,
}

func seedProject(_ *cobra.Command, args []string) {
	program := "go"
	commands := []string{"run", "main.go", "seed:run"}
	// execute command from golang
	if len(args) == 0 || (len(args) == 1 && args[0] == "all") {
		commands = append(commands, "--all")
	} else {
		for _, arg := range args {
			commands = append(commands, "--name")
			commands = append(commands, arg)
		}
	}
	err := utility.ExecuteCommand(program, commands, args...)
	if err != nil {
		return
	}
}
