package cmd

import (
	"github.com/spf13/cobra"
)

// Root is the root of the command execution.
// root can be thought of as a main program.
var Root = &cobra.Command{
	Use:   "geng",
	Short: "A generator for Cobra based Applications",
	Long:  `geng is a CLI library for Go that empowers applications.`,
}

func init() {
	Root.AddCommand(
		newModuleCmd,
		projectCmd,
		runProjectCmd,
		addInfrastructureCmd,
		addServiceCmd,
		seedProjectCmd,
		startProjectCmd,
		migrationProjectCmd,
	)
}
