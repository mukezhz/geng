package main

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"github.com/spf13/cobra"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/wesionary/*
var templatesFS embed.FS

type ModuleData struct {
	ModuleName        string
	PackageName       string
	ProjectModuleName string
	ProjectName       string
	GoVersion         string
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
	newProjectCmd.Flags().StringP("dir", "d", "", "target directory")
	newProjectCmd.Flags().StringP("version", "v", "", "version support")
	rootCmd.AddCommand(newModuleCmd, newProjectCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getModuleDataFromModuleName(moduleName, projectModuleName, goVersion string) ModuleData {
	c := cases.Title(language.English)
	titleModuleName := c.String(moduleName)
	splitedModule := strings.Split(titleModuleName, " ")
	lowerModuleName := ""
	acutalModuleName := ""
	if len(splitedModule) > 0 {
		acutalModuleName = strings.Join(splitedModule, "")
		lowerModuleName = strings.Join(splitedModule, "_")
	} else {
		acutalModuleName = titleModuleName
	}
	lowerModuleName = strings.ToLower(lowerModuleName)
	data := ModuleData{
		ModuleName:        acutalModuleName,
		PackageName:       lowerModuleName,
		ProjectModuleName: projectModuleName,
		ProjectName:       titleModuleName,
		GoVersion:         goVersion,
	}
	return data
}

func getModuleNameFromGoModFile() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			// Extract module name
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	abs, _ := filepath.Abs("go.mod")
	return "", fmt.Errorf("module directive not found in %s", abs)
}

func createModule(cmd *cobra.Command, args []string) {
	projectModuleName, err := getModuleNameFromGoModFile()
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	projectPath, err := findGitRoot(currentDir)
	if err != nil {
		fmt.Println("Error finding Git root:", err)
		return
	}
	mainModulePath := filepath.Join(projectPath, "domain", "features", "module.go")
	if err != nil {
		panic(err)
	}

	moduleName := args[1]
	data := getModuleDataFromModuleName(moduleName, projectModuleName, "")

	// Define the directory structure
	targetRoot := filepath.Join(".", "domain", "features", data.PackageName)
	templatePath := filepath.Join(".", "templates", "wesionary", "module")

	generateFiles(templatePath, targetRoot, data)

	updatedCode := AddAnotherFxOptionsInModule(mainModulePath, data.PackageName, data.ProjectModuleName)
	writeContentToPath(mainModulePath, updatedCode)

}

func createProject(cmd *cobra.Command, args []string) {
	projectName := args[0]
	projectModuleName, _ := cmd.Flags().GetString("mod")
	if projectModuleName == "" {
		panic("project module name is required")
	}
	goVersion, _ := cmd.Flags().GetString("version")
	goVersion = checkVersion(goVersion)
	data := getModuleDataFromModuleName(projectName, projectModuleName, goVersion)
	targetedDirectory, _ := cmd.Flags().GetString("dir")
	if targetedDirectory == "" {
		targetedDirectory = filepath.Join(targetedDirectory, data.PackageName)
	}
	log.Printf("Project Name: %s, Project Module Name: %s", projectName, projectModuleName)
	//root := "./templates/wesionary/project"
	targetRoot := targetedDirectory

	templatePath := "templates/wesionary/project"
	generateFiles(templatePath, targetRoot, data)
	return
}

func generateFiles(templatePath string, targetRoot string, data ModuleData) error {
	return fs.WalkDir(templatesFS, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Create the same structure in the target directory
		relPath, _ := filepath.Rel(templatePath, path)
		targetPath := filepath.Join(targetRoot, filepath.Dir(relPath))
		fileName := filepath.Base(path)
		dst := filepath.Join(targetPath, fileName)
		os.MkdirAll(targetPath, os.ModePerm)

		// Handle template files
		if filepath.Ext(path) == ".tmpl" {
			// Generate the Go file in the target directory
			goFile := strings.Replace(dst, "tmpl", "go", 1)
			generateFromEmbeddedTemplate(path, goFile, data)
		} else {
			// Copy or process other files as before
			if filepath.Ext(path) == ".mod" || filepath.Ext(path) == ".md" {
				dst = strings.Replace(dst, ".mod", "", 1)
				generateFromEmbeddedTemplate(path, dst, data)
			} else {
				if strings.HasPrefix(fileName, "hidden.") {
					dst = strings.Replace(dst, "hidden.", ".", 1)
				}
				// just copy the files to the target directory
				if err := copyFile(path, dst); err != nil {
					panic(err)
				}
			}
		}
		return nil
	})
}

func generateFromEmbeddedTemplate(path, targetFilePath string, data interface{}) {
	tmpl, err := template.ParseFS(templatesFS, path)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(targetFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := templatesFS.Open(src)
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

func checkVersion(goVersion string) string {
	if goVersion == "" {
		return goVersion
	}
	split := strings.Split(goVersion, ".")
	if len(split) >= 3 {
		return strings.Join(split[:2], ".")
	}
	return goVersion
}

func AddAnotherFxOptionsInModule(path, module, projectModule string) string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	importPackage(node, projectModule, module)

	// Traverse the AST and find the fx.Options call
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Module" {
					x.Args = append(x.Args, &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("fx"),
							Sel: ast.NewIdent("Options"),
						},
						Args: []ast.Expr{
							ast.NewIdent(module + ".Module"),
						},
					})
				}
			}
		}
		return true
	})

	// Add the source code in buffer
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		log.Fatal(err)
	}
	formattedCode := buf.String()
	providerToInsert := fmt.Sprintf("fx.Options(%v.Module),", module)
	formattedCode = strings.Replace(formattedCode, providerToInsert, "\n\t"+providerToInsert, 1)
	return formattedCode
}

func findGitRoot(path string) (string, error) {
	if path == "/" {
		return "", fmt.Errorf("reached root directory, .git not found")
	}

	// Check if .git directory exists in the current path
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		return path, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	parentDir := filepath.Dir(path)
	return findGitRoot(parentDir)
}

func writeContentToPath(path, content string) {
	err := os.WriteFile(path, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func importPackage(node *ast.File, projectModule, packageName string) {
	importSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%v/domain/features/%v"`, projectModule, packageName),
		},
	}

	importDecl := &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: token.Pos(1), // for grouping
		Specs:  []ast.Spec{importSpec},
	}

	// Check if there are existing imports, and if so, add to them
	found := false
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.IMPORT {
			genDecl.Specs = append(genDecl.Specs, importSpec)
			found = true
			break
		}
	}

	// If no import declaration exists, add the new one to Decls
	if !found {
		node.Decls = append([]ast.Decl{importDecl}, node.Decls...)
	}
}
