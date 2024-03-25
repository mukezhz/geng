package cmd

func init() {
	setupFlagsForNewProject(newProjectCmd)
	rootCmd.AddCommand(newModuleCmd)
	rootCmd.AddCommand(newProjectCmd)
	rootCmd.AddCommand(runProjectCmd)
	rootCmd.AddCommand(addInfrastructureCmd)
	rootCmd.AddCommand(addServiceCmd)
	rootCmd.AddCommand(seedProjectCmd)
	rootCmd.AddCommand(startProjectCmd)
	rootCmd.AddCommand(migrationProjectCmd)
}
