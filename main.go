package main

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/terminal"
	"github.com/spf13/cobra"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

const (
	ProjectNameKEY        = "projectName"
	ProjectModuleNameKEY  = "projectModuleName"
	AuthorKEY             = "author"
	ProjectDescriptionKEY = "projectDescription"
	GoVersionKEY          = "goVersion"
	DirectoryKEY          = "directory"
)

const (
	ProjectName        = "Project Name"
	ProjectModuleName  = "Project Module"
	Author             = "Author Detail"
	ProjectDescription = "Project Description"
	GoVersion          = "Go Version"
	Directory          = "Project Directory"
)

//go:embed templates/wesionary/*
var templatesFS embed.FS

type ModuleData struct {
	ModuleName         string
	PackageName        string
	ProjectModuleName  string
	ProjectName        string
	GoVersion          string
	ProjectDescription string
	Author             string
	Directory          string
}

type GoMod struct {
	Module    string
	GoVersion string
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
	Args:  cobra.MaximumNArgs(1),
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

func getModuleNameFromGoModFile() (GoMod, error) {
	file, err := os.Open("go.mod")
	goMod := GoMod{}

	if err != nil {
		return goMod, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			// Extract module name
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				goMod.Module = parts[1]
			}
		} else if strings.HasPrefix(line, "go ") {
			// Extract Go version
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				goMod.GoVersion = parts[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return goMod, err
	}

	abs, _ := filepath.Abs("go.mod")
	return goMod, fmt.Errorf("module directive not found in %s", abs)
}

func createModule(_ *cobra.Command, args []string) {
	projectModule, err := getModuleNameFromGoModFile()
	currentDir, err := os.Getwd()
	if err != nil {
		color.Redln("Error getting current directory:", err)
		panic(err)
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
	if !checkGolangIdentifier(moduleName) {
		color.Redln("Error: module name is invalid")
		return
	}
	data := getModuleDataFromModuleName(moduleName, projectModule.Module, projectModule.GoVersion)

	// Define the directory structure
	targetRoot := filepath.Join(".", "domain", "features", data.PackageName)
	templatePath := filepath.Join(".", "templates", "wesionary", "module")

	err = generateFiles(templatePath, targetRoot, data)
	if err != nil {
		color.Redln("Error: generate file", err)
		return
	}

	updatedCode := AddAnotherFxOptionsInModule(mainModulePath, data.PackageName, data.ProjectModuleName)
	writeContentToPath(mainModulePath, updatedCode)

	PrintColorizeModuleDetail(data)

}

func createProject(cmd *cobra.Command, args []string) {
	var projectName string
	var projectModuleName string
	var goVersion string
	var projectDescription string
	var author string
	var directory string

	if len(args) == 0 {
		questions := []terminal.ProjectQuestion{
			terminal.NewShortQuestion(ProjectNameKEY, ProjectName+" *", "Enter Project Name:"),
			terminal.NewShortQuestion(ProjectModuleNameKEY, ProjectModuleName+" *", "Enter Module Name:"),
			terminal.NewShortQuestion(AuthorKEY, Author+" [Optional]", "Enter Author Detail[Mukesh Chaudhary <mukezhz@duck.com>] [Optional]"),
			terminal.NewLongQuestion(ProjectDescriptionKEY, ProjectDescription+" [Optional]", "Enter Project Description [Optional]"),
			terminal.NewShortQuestion(GoVersionKEY, GoVersion+" [Optional]", "Enter Go Version (Default: 1.20) [Optional]"),
			terminal.NewShortQuestion(DirectoryKEY, Directory+" [Optional]", "Enter Project Directory (Default: package_name) [Optional]"),
		}
		terminal.StartInteractiveTerminal(questions)

		for _, q := range questions {
			switch q.Key {
			case ProjectNameKEY:
				projectName = q.Answer
				break
			case ProjectDescriptionKEY:
				projectDescription = q.Answer
				break
			case AuthorKEY:
				author = q.Answer
				break
			case ProjectModuleNameKEY:
				projectModuleName = q.Answer
				break
			case GoVersionKEY:
				goVersion = q.Answer
				break
			case DirectoryKEY:
				directory = q.Answer
				break
			}
		}
	} else {
		projectName = args[0]
		projectModuleName, _ = cmd.Flags().GetString("mod")
		goVersion, _ = cmd.Flags().GetString("version")
		directory, _ = cmd.Flags().GetString("dir")
	}

	goVersion = checkVersion(goVersion)
	if projectName == "" {
		color.Redln("Error: project name is required")
		return
	}
	if projectModuleName == "" {
		color.Redln("Error: module name is required")
		return
	}

	data := getModuleDataFromModuleName(projectName, projectModuleName, goVersion)
	data.ProjectDescription = projectDescription
	data.Author = author

	data.Directory = directory
	if data.Directory == "" {
		data.Directory = filepath.Join(data.Directory, data.PackageName)
	}
	targetRoot := data.Directory

	templatePath := filepath.Join("templates", "wesionary", "project")
	err := generateFiles(templatePath, targetRoot, data)
	if err != nil {
		color.Redln("Error generate file", err)
		return
	}

	PrintColorizeProjectDetail(data)
	fmt.Println("\n")
	return
}

func PrintColorizeProjectDetail(data ModuleData) {
	color.Cyanln(`
	    GENG: GENERATE GOLANG PROJECT
	
	 ██████╗ ███████╗███╗   ██╗       ██████╗ 
	██╔════╝ ██╔════╝████╗  ██║      ██╔════╝ 
	██║  ███╗█████╗  ██╔██╗ ██║█████╗██║  ███╗
	██║   ██║██╔══╝  ██║╚██╗██║╚════╝██║   ██║
	╚██████╔╝███████╗██║ ╚████║      ╚██████╔╝
	 ╚═════╝ ╚══════╝╚═╝  ╚═══╝       ╚═════╝ 
											  
	
	`)
	color.Greenln("\tThe information you have provided:\n")
	color.Cyanf("\t%-20s💻: %-15s\n", ProjectName, data.ProjectName)
	color.Cyanf("\t%-20s📂: %-15s\n", ProjectModuleName, data.ProjectModuleName)
	color.Cyanf("\t%-20s📚: %-15s\n", ProjectDescription, data.ProjectDescription)
	color.Cyanf("\t%-20s🆚: %-15s\n", GoVersion, data.GoVersion)
	color.Cyanf("\t%-20s🤓: %-15s\n", Author, data.Author)
	PrintFinalStepAfterProjectInitialization(data)
	color.Redln("\n\tThank You For using 🙏🇳🇵🙏:\n")

}
func PrintColorizeModuleDetail(data ModuleData) {
	color.Cyanln(`
	    GENG: GENERATE GOLANG MODULE
	
	 ██████╗ ███████╗███╗   ██╗       ██████╗ 
	██╔════╝ ██╔════╝████╗  ██║      ██╔════╝ 
	██║  ███╗█████╗  ██╔██╗ ██║█████╗██║  ███╗
	██║   ██║██╔══╝  ██║╚██╗██║╚════╝██║   ██║
	╚██████╔╝███████╗██║ ╚████║      ╚██████╔╝
	 ╚═════╝ ╚══════╝╚═╝  ╚═══╝       ╚═════╝ 
											  
	
	`)
	color.Greenln("\tThe information you have provided:\n")
	color.Cyanf("\t%-20s💻: %-15s\n", ProjectName, data.ProjectName)
	color.Cyanf("\t%-20s📂: %-15s\n", ProjectModuleName, data.ProjectModuleName)
	color.Cyanf("\t%-20s🆚: %-15s\n", GoVersion, data.GoVersion)
	PrintFinalStepAfterModuleInitialization(data)
	color.Redln("\n\tThank You For using 🙏🇳🇵🙏:\n")
}

func PrintFinalStepAfterModuleInitialization(data ModuleData) {
	output := fmt.Sprintf(`
	🎉 Successfully created module %v

	↪️ Restart the server to see the changes:

	🌐 Navigate to the following path:
	    %v

`, data.ModuleName, "/api/"+data.PackageName)
	color.Yellowf(output)
}

func PrintFinalStepAfterProjectInitialization(data ModuleData) {
	output := fmt.Sprintf(`
	💻 Change directory to project:
	    cd %v
	
	💾 Initalize git repository:
	    git init

	📚 Sync dependencies:
	    go mod tidy
	
	🕵 Copy .env.example to .env:
	    cp .env.example .env
	
	🏃 Start Project 🏃:
	    go run main.go app:serve
`, data.Directory)
	color.Yellowf(output)
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
		err = os.MkdirAll(targetPath, os.ModePerm)
		if err != nil {
			color.Redln("Error make dir", err)
			return err
		}

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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			color.Println("Error closing", err)
		}
	}(file)

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
	defer func(sourceFile fs.File) {
		err := sourceFile.Close()
		if err != nil {
			color.Redln("Error closing:", err)
		}
	}(sourceFile)

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			color.Redln("Error closing:", err)
		}
	}(destFile)

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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			color.Redln("Error closing:", err)
		}
	}(file)

	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
}

func checkVersion(goVersion string) string {
	versionRegex := regexp.MustCompile(`^\d+\.\d+(\.\d+)?$`)
	if goVersion == "" {
		return "1.20"
	}
	if !versionRegex.MatchString(goVersion) {
		return "error"
	}
	split := strings.Split(goVersion, ".")
	if len(split) >= 3 {
		return strings.Join(split[:2], ".")
	}
	return goVersion
}

func checkGolangIdentifier(identifier string) bool {
	if identifier == "" {
		return false
	}

	for i, r := range identifier {
		if i == 0 && !unicode.IsLetter(r) && r != '_' {
			return false
		}
		if i > 0 && !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}

	return true
}

func AddAnotherFxOptionsInModule(path, module, projectModule string) string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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
	path := filepath.Join(projectModule, "domain", "features", packageName)
	importSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%v"`, path),
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
