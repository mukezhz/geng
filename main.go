package main

import (
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
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

func getModuleDataFromModuleName(moduleName string) ModuleData {
	lowerModuleName := strings.ToLower(moduleName)
	c := cases.Title(language.English)
	titleModuleName := c.String(moduleName)
	data := ModuleData{
		ModuleName:        titleModuleName,
		PackageName:       lowerModuleName,
		ProjectModuleName: "github.com/mukezhz/geng",
	}
	return data
}

func createModule(cmd *cobra.Command, args []string) {
	moduleName := args[1]
	data := getModuleDataFromModuleName(moduleName)

	// Define the directory structure
	baseDir := filepath.Join(".", "domain", data.ModuleName)
	templateDir := filepath.Join(".", "templates", "wesionary", "domain", "features")
	dir, err := os.ReadDir("domain")
	if err != nil {
		// If the directory does not exist, ignore the error
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
	if err := os.Mkdir(filepath.Join(".", "domain", data.PackageName), 0755); err != nil {
		panic(err)
	}
	readDir, err := os.ReadDir(templateDir)
	if err != nil {
		panic(err)
	}

	for _, file := range readDir {
		log.Printf("%#v\n\n", file)
		goFile := strings.Replace(file.Name(), "tmpl", "go", 1)
		generateFromTemplate(filepath.Join(templateDir, file.Name()), filepath.Join(baseDir, goFile), data)
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
