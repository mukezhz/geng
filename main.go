package main

import (
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ModuleData struct {
	ModuleName        string
	PackageName       string
	ProjectModuleName string
}

var rootCmd = &cobra.Command{
	Use:   "geng",
	Short: "Go Generate[geng] is a tool for generating Go modules",
}

var newModuleCmd = &cobra.Command{
	Use:   "gen module [name]",
	Short: "Create a new module",
	Args:  cobra.ExactArgs(2),
	Run:   createModule,
}

func init() {
	rootCmd.AddCommand(newModuleCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func createModule(cmd *cobra.Command, args []string) {
	moduleName := args[1]
	lowerModuleName := strings.ToLower(moduleName)
	c := cases.Title(language.English)
	titleModuleName := c.String(moduleName)
	data := ModuleData{
		ModuleName:        titleModuleName,
		PackageName:       lowerModuleName,
		ProjectModuleName: "github.com/mukezhz/geng",
	}

	// Define the directory structure
	baseDir := filepath.Join(".", "domain", lowerModuleName)
	templateDir := filepath.Join(".", "templates", "wesionary")
	dir, err := os.ReadDir("domain")
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			panic(err)
		}
	}
	if len(dir) == 0 {
		if err := os.Mkdir(filepath.Join(".", "domain"), 0755); err != nil {
			panic(err)
		}
	}
	// Create directories
	if err := os.Mkdir(filepath.Join(".", "domain", lowerModuleName), 0755); err != nil {
		panic(err)
	}

	moduleTemplates := []string{"controller.tmpl", "service.tmpl", "model.tmpl", "route.tmpl", "module.tmpl"}
	goFiles := []string{"controller.go", "service.go", "model.go", "route.go", "module.go"}
	for n, tmpl := range moduleTemplates {
		generateFromTemplate(filepath.Join(templateDir, tmpl), filepath.Join(baseDir, goFiles[n]), data)
	}
}

func generateFromTemplate(templateFile, outputFile string, data ModuleData) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
}
