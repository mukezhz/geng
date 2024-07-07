package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/mukezhz/bru-go/bru"
	"github.com/mukezhz/bru-go/parser"
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
			tokens := bru.GetTokens(string(content))
			ast := bru.GetAST(tokens)
			m := utility.GetModuleNameFromPath(path)
			httpNode := bru.GetTagNode[parser.HTTPNode](ast)
			metaNode := bru.GetTagNode[parser.MetaNode](ast)
			b := model.BruModel{
				Route:       utility.SanitizeEndpoint(httpNode.URL),
				Method:      strings.ToUpper(httpNode.Method),
				Body:        httpNode.Body,
				Handler:     utility.ToPasalCase(metaNode.Name),
				ModuleName:  m,
				Name:        metaNode.Name,
				Description: metaNode.Name,
			}

			if !utility.FileExists("domain/" + m) {

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

		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", rootDir, err)
	}
}
