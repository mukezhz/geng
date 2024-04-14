package cmd

import (
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var generateFxCmd = &cobra.Command{
	Use:   "fx",
	Short: "Generate the fx configuration file for the project",
	Long: `
Generate a module by reading comments from the source code.
Example: 
  geng fx

	`,
	Args: cobra.MaximumNArgs(2),
	Run:  generateFx,
}

func generateFx(_ *cobra.Command, _ []string) {
	utility.GenerateFxModule()
}
