package cmd

import (
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", false, "use Viper for configuration")
	cobra.CheckErr(viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author")))
	cobra.CheckErr(viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper")))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "none")
	newProjectCmd.Flags().StringP("mod", "m", "", "features name")
	newProjectCmd.Flags().StringP("dir", "d", "", "target directory")
	newProjectCmd.Flags().StringP("version", "v", "", "version support")
	rootCmd.AddCommand(newModuleCmd)
	rootCmd.AddCommand(newProjectCmd)
	rootCmd.AddCommand(runProjectCmd)
}
