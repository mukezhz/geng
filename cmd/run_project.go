package cmd

import (
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var runProjectCmd = &cobra.Command{
	Use:   "run [project name]",
	Short: "Run the project",
	Args:  cobra.MaximumNArgs(1),
	Run:   runProject,
}

func runProject(_ *cobra.Command, args []string) {
	runGo := "go"
	// execute command from golang
	err := utility.ExecuteCommand(runGo, args...)
	if err != nil {
		return
	}
}

func init() {
	newProjectCmd.Flags().StringP("mod", "m", "", "features name")
	newProjectCmd.Flags().StringP("dir", "d", "", "target directory")
	newProjectCmd.Flags().StringP("version", "v", "", "version support")
	rootCmd.AddCommand(newModuleCmd)
	rootCmd.AddCommand(newProjectCmd)
	rootCmd.AddCommand(runProjectCmd)
	rootCmd.AddCommand(addInfrastructureCmd)
	rootCmd.AddCommand(addServiceCmd)
}
