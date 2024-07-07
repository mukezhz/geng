package cmd

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/model"
	"github.com/mukezhz/geng/pkg/utility"
)

func getProjectPath() (string, *model.GoMod) {
	projectModule, err := utility.GetModuleNameFromGoModFile()
	if err != nil {
		fmt.Println("Error finding Module name from go.mod:", err)
		return "", nil
	}
	currentDir, err := os.Getwd()
	if err != nil {
		color.Redln("Error getting current directory:", err)
		panic(err)
	}
	projectPath, err := utility.FindGitRoot(currentDir)
	if err != nil {
		fmt.Println("Error finding Git root:", err)
		return "", &projectModule
	}
	return projectPath, &projectModule
}
