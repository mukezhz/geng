package utility

import (
	"embed"
	"log"
	"os"
	"text/template"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/model"
)

func GenerateFromTemplate(templateFile, outputFile string, data model.ModuleData) {
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

func ListDirectory(templatesFS embed.FS, dirPath string) []string {
	f := []string{}
	files, err := templatesFS.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		f = append(f, file.Name())
	}
	return f
}
