package cmd

import (
	"embed"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "geng",
		Short: "A generator for Cobra based Applications",
		Long:  `geng is a CLI library for Go that empowers applications.`,
	}
	templatesFS embed.FS
)

func Execute(fs embed.FS) error {
	templatesFS = fs
	return rootCmd.Execute()
}
