package gen

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/model"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/mukezhz/geng/templates"
)

// ProjectGenerator is argument for project generator
type ProjectGenerator struct {
	Name        string
	ModuleName  string
	Description string
	Author      string
	Directory   string
	GoVersion   string
}

// Fill fills up items from map to the struct
func (p *ProjectGenerator) Fill(d map[string]string) {
	p.Name, _ = d[constant.ProjectNameKEY]
	p.ModuleName, _ = d[constant.ProjectModuleNameKEY]
	p.Author, _ = d[constant.AuthorKEY]
	p.Description, _ = d[constant.ProjectDescriptionKEY]
	p.GoVersion, _ = d[constant.GoVersionKEY]
	p.Directory, _ = d[constant.DirectoryKEY]
}

// Validate validates generated project arguments
func (p *ProjectGenerator) Validate() error {
	p.GoVersion = utility.CheckVersion(p.GoVersion)
	if p.Name == "" {
		return errors.New("Error: project name is required")
	}
	if p.ModuleName == "" {
		return errors.New("Error: golang module name is required")
	}
	return nil
}

// Generate genrates the project given the arguments in the struct
func (p *ProjectGenerator) Generate() (*model.ModuleData, error) {
	data := utility.GetModuleDataFromModuleName(p.Name, p.ModuleName, p.GoVersion)
	data.ProjectDescription = p.Description
	data.Author = p.Author

	data.Directory = p.Directory
	if data.Directory == "" {
		data.Directory = filepath.Join(data.Directory, data.PackageName)
	}

	targetRoot := data.Directory

	templatePath := utility.IgnoreWindowsPath(filepath.Join("templates", "wesionary", "project"))
	err := utility.GenerateFiles(templates.FS, templatePath, targetRoot, data)
	if err != nil {
		return nil, fmt.Errorf("Error generating file: %v", err)
	}

	utility.PrintColorizeProjectDetail(data)
	fmt.Println("")

	return &data, nil
}
