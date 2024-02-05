package cmd

import (
	"embed"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
	templatesFS embed.FS
)

func Execute(fs embed.FS) error {
	templatesFS = fs
	return rootCmd.Execute()
}
