package gen

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/mukezhz/geng/pkg/constant"
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

	Infra *InfraGenerator
}

// Fill fills up items from map to the struct
func (p *ProjectGenerator) Fill(d map[string]string) {
	p.Name, _ = d[constant.ProjectNameKEY]
	p.ModuleName, _ = d[constant.ProjectModuleNameKEY]
	p.Author, _ = d[constant.AuthorKEY]
	p.Description, _ = d[constant.ProjectDescriptionKEY]
	p.GoVersion, _ = d[constant.GoVersionKEY]
	p.Directory, _ = d[constant.DirectoryKEY]

	p.Infra = &InfraGenerator{
		Directory: p.Directory,
	}
}

// FillProjectMetadataFromJson fills up project metadata from json file
func (p *ProjectGenerator) FillProjectMetadataFromJson() {
	data := utility.ReadJsonFile("geng.json")
	for k, v := range data {
		switch k {
		case constant.ProjectNameKEY:
			p.Name = v.(string)
		case constant.ProjectModuleNameKEY:
			p.ModuleName = v.(string)
		case constant.AuthorKEY:
			p.Author = v.(string)
		case constant.ProjectDescriptionKEY:
			p.Description = v.(string)
		case constant.GoVersionKEY:
			p.GoVersion = v.(string)
		case constant.DirectoryKEY:
			p.Directory = v.(string)
		}
	}
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
func (p *ProjectGenerator) Generate(selectedInfra []int) error {
	data := utility.GetModuleDataFromModuleName(p.Name, p.ModuleName, p.GoVersion)
	data.ProjectDescription = p.Description
	data.Author = p.Author

	data.Directory = p.Directory
	if data.Directory == "" {
		data.Directory = filepath.Join(data.Directory, data.PackageName)
	}

	p.Infra.Directory = data.Directory

	templatePath := utility.IgnoreWindowsPath(filepath.Join("templates", "wesionary", "project"))
	err := utility.GenerateFiles(templates.FS, templatePath, p.Infra.Directory, data)
	if err != nil {
		return fmt.Errorf("error generating file: %v", err)
	}

	utility.PrintColorizeProjectDetail(data)
	fmt.Println("")

	if err := p.Infra.Validate(); err != nil {
		return err
	}

	if len(selectedInfra) == 0 {
		return errors.New("no infrastructure selected")
	}

	selectedInfras := p.Infra.GetSelectedItems(selectedInfra)
	if err := p.Infra.Generate(data, selectedInfra); err != nil {
		return fmt.Errorf("generation error: %v", err)
	}

	utility.PrintColorizeInfrastructureDetail(data, selectedInfras)
	return nil
}
