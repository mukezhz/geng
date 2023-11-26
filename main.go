package main

import (
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
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
	Use:   "gen features [name]",
	Short: "Create a new features",
	Args:  cobra.ExactArgs(2),
	Run:   createModule,
}

var newProjectCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	Run:   createProject,
}

func init() {
	newProjectCmd.Flags().StringP("mod", "m", "", "features name")
	rootCmd.AddCommand(newModuleCmd, newProjectCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getModuleDataFromModuleName(moduleName, projectModuleName string) ModuleData {
	lowerModuleName := strings.ToLower(moduleName)
	c := cases.Title(language.English)
	titleModuleName := c.String(moduleName)
	data := ModuleData{
		ModuleName:        titleModuleName,
		PackageName:       lowerModuleName,
		ProjectModuleName: projectModuleName,
	}
	return data
}

func createModule(cmd *cobra.Command, args []string) {
	moduleName := args[1]
	data := getModuleDataFromModuleName(moduleName, "github.com/mukezhz/test")

	// Define the directory structure
	baseDir := filepath.Join(".", "domain", data.ModuleName)
	templateDir := filepath.Join(".", "templates", "wesionary", "module")
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
		goFile := strings.Replace(file.Name(), "tmpl", "go", 1)
		generateFromTemplate(filepath.Join(templateDir, file.Name()), filepath.Join(baseDir, goFile), data)
	}
}

func createProject(cmd *cobra.Command, args []string) {
	projectName := args[0]
	projectModuleName, _ := cmd.Flags().GetString("mod")
	log.Printf("projectName: %s, projectModulename: %s", projectName, projectModuleName)
	root := "./templates/wesionary/project"
	targetRoot := "generated"
	data := getModuleDataFromModuleName(projectName, projectModuleName)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".tmpl" {
			// Create the same structure in the target directory
			relPath, _ := filepath.Rel(root, path)
			targetPath := filepath.Join(targetRoot, filepath.Dir(relPath))

			// Create the directory if it does not exist
			os.MkdirAll(targetPath, os.ModePerm)

			// Generate the Go file in the target directory
			goFile := strings.Replace(filepath.Join(targetPath, filepath.Base(path)), "tmpl", "go", 1)
			generateFromTemplate(path, goFile, data)
		} else {
			// just copy the files to the target directory
			relPath, _ := filepath.Rel(root, path)
			targetPath := filepath.Join(targetRoot, relPath)

			os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)

			// Copy the file
			if err := copyFile(path, targetPath); err != nil {
				panic(err)
			}
		}
		return nil
	})
	return

}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
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
