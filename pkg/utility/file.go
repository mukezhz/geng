package utility

import (
	"bufio"
	"embed"
	"fmt"
	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func WriteContentToPath(path, content string) {
	err := os.WriteFile(path, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func FindGitRoot(path string) (string, error) {
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
	return FindGitRoot(parentDir)
}

func GetModuleDataFromModuleName(moduleName, projectModuleName, goVersion string) model.ModuleData {
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
	data := model.ModuleData{
		ModuleName:        acutalModuleName,
		PackageName:       lowerModuleName,
		ProjectModuleName: projectModuleName,
		ProjectName:       titleModuleName,
		GoVersion:         goVersion,
	}
	return data
}

func GenerateFiles(templatesFS embed.FS, templatePath string, targetRoot string, data model.ModuleData) error {
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
			generateFromEmbeddedTemplate(templatesFS, path, goFile, data)
		} else {
			// Copy or process other files as before
			if filepath.Ext(path) == ".mod" || filepath.Ext(path) == ".md" {
				dst = strings.Replace(dst, ".mod", "", 1)
				generateFromEmbeddedTemplate(templatesFS, path, dst, data)
			} else {
				if strings.HasPrefix(fileName, "hidden.") {
					dst = strings.Replace(dst, "hidden.", ".", 1)
				}
				// just copy the files to the target directory
				if err := copyFile(path, dst, templatesFS); err != nil {
					panic(err)
				}
			}
		}
		return nil
	})
}

func generateFromEmbeddedTemplate(templatesFS embed.FS, path, targetFilePath string, data interface{}) {
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

func copyFile(src, dst string, templatesFS embed.FS) error {
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

func CheckVersion(goVersion string) string {
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

func GetModuleNameFromGoModFile() (model.GoMod, error) {
	file, err := os.Open("go.mod")
	goMod := model.GoMod{}

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
