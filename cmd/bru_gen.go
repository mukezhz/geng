package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/mukezhz/bru-go/bru"
	"github.com/mukezhz/geng/pkg/gen"
	"github.com/mukezhz/geng/pkg/model"
	"github.com/mukezhz/geng/pkg/utility"
	"github.com/spf13/cobra"
)

var bruCmd = &cobra.Command{
	Use:   "bru [path]",
	Short: "Generate from Bru file",
	Args:  cobra.MaximumNArgs(2),
	Run:   generateFromBru,
}

func generateFromBru(_ *cobra.Command, args []string) {
	WalkBruFiles(args[0])
}

func WalkBruFiles(rootDir string) {

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".bru") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			bru, err := bru.Unmarshal(content)
			if err != nil {
				return err
			}

			m := utility.GetModuleNameFromPath(path)
			b := model.BruModel{
				Route:       utility.SanitizeEndpoint(bru.HTTP.URL),
				Method:      strings.ToUpper(bru.HTTP.Method),
				Body:        bru.Body.Content,
				Handler:     utility.ToPasalCase(bru.Meta.Name),
				ModuleName:  m,
				Name:        bru.Meta.Name,
				Description: bru.Meta.Name,
			}

			if !utility.FileExists(filepath.Join("domain", m)) {
				projectPath, projectModule := getProjectPath()
				if projectModule == nil || projectPath == "" {
					return nil
				}

				if !utility.FileExists(filepath.Join(projectPath, "domain", "module.go")) {
					color.Redln("[Project not found...]")
					return errors.New("please initialize the project first")
				}
				mainModulePath := filepath.Join(projectPath, "domain", "module.go")
				if !utility.FileExists(mainModulePath) {
					color.Redln("Error: module.go not found")
					return nil
				}
				// create a module
				generateModule(mainModulePath, []string{"", m}, *projectModule)
			}

			gen.AddRoute(b)
			gen.AddController(b)
			utility.PrintGenerationFromBrufile()

		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", rootDir, err)
	}
}
